<?php

// ==================================== Ingestão da request que vem do stdin =============

$_http_header = fgets(STDIN);
$_http_header = trim($_http_header);
$_http_header = explode(" ", $_http_header);

$_SERVER['REQUEST_METHOD'] = $_http_header[0];
$_SERVER['REQUEST_URI'] = $_http_header[1];
$_SERVER['QUERY_STRING'] = parse_url($_http_header[1], PHP_URL_QUERY);

$_HEADERS = array();
function header(string $header) {
    global $_HEADERS;
    array_push($_HEADERS, $header);
}

while (true) {
    $_header = fgets(STDIN);
    $_header = trim($_header);
    if (strlen($_header) == 0) {
        break;
    }
    $_header_name = explode(":", $_header)[0];
    $_header_value = substr($_header, strlen($_header_name)+1);
    $_header_value = trim($_header_value);

    // fixes security issue where an attacker could
    // pass arbitrary stuff into the TAILSCALE_USER_LOGIN header
    // error_log("header: $_header_name");
    if ($_header_name == "Tailscale-User-Login") {
        define("TS_LOGIN", $_header_value);
    }

    $_header_name = strtoupper($_header_name);
    $_header_name = str_replace("-", "_", $_header_name);

    $_SERVER['HTTP_' . $_header_name] = $_header_value;
}

// ==================================== Primitiva de header pra retorno =============
set_header("Server", "GambiarraStation");
set_header("Connection", "close");

$_HEADERS_KV = array();
function set_header(string $key, string $value) {
    global $_HEADERS_KV;
    $_HEADERS_KV[$key] = $value;
}

// ==================================== Utilitários para content-type =============
function set_contenttype(string $content_type) {
    set_header("Content-Type", $content_type);
}

set_contenttype("auto"); // default

function content_text() {
    set_contenttype("text/plain; charset=utf-8");
}

function content_html() {
    set_contenttype("text/html");
}

function mime_from_buffer($buffer) {
    $finfo = new finfo(FILEINFO_MIME_TYPE);
    return $finfo->buffer($buffer);
}


// ==================================== Utilitários para roteamento =============
$INPUT_DATA = array_merge_recursive($_GET, $_POST);
$ROUTE = parse_url($_SERVER["REQUEST_URI"])["path"];
$IS_ROUTED = false;

function mark_routed() {
    global $IS_ROUTED;
    $orig_VALUE = $IS_ROUTED;
    $IS_ROUTED = true;
    return $orig_VALUE;
}


/**
 * bring required variables to scope then require the new file
 */
function execphp(string $script) {
    global $INPUT_DATA, $ROUTE;
    require $script;
}

/**
 * create a route that matches anything starting with $base_route
 */
function use_route(string $base_route, string $handler_script) {
    global $ROUTE, $IS_ROUTED;
    if (str_starts_with($ROUTE, $base_route)) {
        if ($IS_ROUTED) return;
        $ROUTE = substr($ROUTE, strlen($base_route));
        execphp($handler_script);
    }
}

/**
 * create a route that matches exactly $selected_route
 */
function exact_route(string $selected_route, string $handler_script) {
    global $ROUTE;
    if (strcmp($ROUTE, $selected_route) == 0) {
        if (mark_routed()) return;
        execphp($handler_script);
    }
}

/**
 * create a route with route params, like /users/:user/info 
 * and pass the route param with input_data
 */
function exact_with_route_param(string $selected_route, string $handler_script) {
    global $INPUT_DATA, $ROUTE;
    $preprocess = function (string $raw_route) {
        $splitted = preg_split("/\//", $raw_route);
        $splitted = array_filter($splitted, function($v, $k) {
            return !is_empty_string($v);
        }, ARRAY_FILTER_USE_BOTH);
        return array_values($splitted);
    };
    $params_parts = $preprocess($selected_route);
    $route_parts = $preprocess($ROUTE);

    $extra_params = [];
    if (count($params_parts) == count($route_parts)) {
        for ($i = 0; $i < count($params_parts); $i++) {
            if (str_starts_with($params_parts[$i], ":")) {
                $extra_params[substr($params_parts[$i], 1)] = $route_parts[$i];
            } else {
                if (strcmp($params_parts[$i], $route_parts[$i]) != 0) {
                    return;
                }
            }
        }
    } else {
        return;
    }
    if (mark_routed()) return;
    $INPUT_DATA = array_merge_recursive($INPUT_DATA, $extra_params);
    execphp($handler_script);
}

function content_scope_push() {
    ob_start();
}

function content_scope_pop() {
    $data = ob_get_contents();
    ob_end_clean();
    return $data;
}

function content_scope_pop_markdown() {
    content_html(); // would be always html anyway
    $lines = content_scope_pop();

    $lines = preg_replace("/\n\#/", "\n\n#", $lines);
    $lines = preg_replace("/\n+/", "\n", $lines);

    $lines = preg_replace('/\*\*(.*)\*\*/', '<b>\\1</b>', $lines);
    $lines = preg_replace('/\_\_(.*)\_\_/', '<b>\\1</b>', $lines);
    $lines = preg_replace('/\*(.*)\*/', '<em>\\1</em>', $lines);
    $lines = preg_replace('/\_(.*)\_/', '<em>\\1</em>', $lines);
    $lines = preg_replace('/\~(.*)\~/', '<del>\\1</del>', $lines);

    while (true) {
        if (!preg_match('/\!\[([^\]]*?)\]\(([a-z]*:\/\/([a-z-0-9]*\.?)+(:[0-9]+)?[^\s\)]*)\)/m', $lines, $link_found, 0)) {
            break;
        }
        $search_term = $link_found[0];
        $label = $link_found[1];
        $link = $link_found[2];
        $json = json_encode($link_found);
        content_scope_push();
        echo '<img alt="';
        echo $label;
        echo '" title="';
        echo $label;
        echo '" src="';
        echo $link;
        echo '">';
        $replace_term = content_scope_pop();
        $lines = str_replace($search_term, $replace_term, $lines);
    }
    while (true) {
        if (!preg_match('/[\(\s]([a-z]*:\/\/([a-z-0-9]*\.?)+(:[0-9]+)?[^\s\)]*)[\)\s]/m', $lines, $link_found, PREG_OFFSET_CAPTURE, 0)) {
            break;
        }
        $link = substr($link_found[0][0], 1, -1);
        $offset = $link_found[0][1] + 1;
        error_log("match: link='$link' offset='$offset'");
        if (substr($lines, $offset - 2, 2) == "](") {
            $exploded_label = explode("[", substr($lines, 0, $offset - 2));
            $label = $exploded_label[array_key_last($exploded_label)];
            $search_term = "[" . $label . "](" . $link . ")";
            content_scope_push();
            echo '<a href="';
            echo $link;
            echo '">';
            echo $label;
            echo "</a>";
            $replace_term = content_scope_pop();
            $lines = str_replace($search_term, $replace_term, $lines);
        } else {
            $search_term = $link;
            content_scope_push();
            echo '<a href="';
            echo $link;
            echo '">';
            echo $link;
            echo '</a>';
            $replace_term = content_scope_pop();
            $lines = str_replace($search_term, $replace_term, $lines);
        }
    }

    content_scope_push(); // output accumulator
    // $replaces = array()

    $lines = explode("\n", $lines);

    foreach ($lines as $i => $line) {
        $line = trim($line);
        if ($line == "") {
            continue;
        }
        $tag = 'p';
        preg_match('/^#*/', $line, $heading_level);
        $heading_level = strlen($heading_level[0]);

        if ($heading_level > 0) {
            $tag = "h$heading_level";
            $line = substr($line, $heading_level);
            $line = trim($line);
        } else {
            preg_match('/^>*/', $line, $blockquote_level);
            $blockquote_level = strlen($blockquote_level[0]);
            $line = substr($line, $blockquote_level);
            $line = trim($line);
            for ($i = 0; $i < $blockquote_level; $i++) {
                $line = "<blockquote>$line</blockquote>";
            }
        }
        $line = "<$tag>$line</$tag>";
        echo $line;
    }
    return content_scope_pop();
}

function auth_tailscale() {
    $name = "";
    $profile_pic = "";
    $has_login = true;
    // if (array_key_exists("HTTP_Tailscale-User-Login", $_SERVER)) {
    //     $login = $_SERVER["HTTP_Tailscale-User-Login"];
    // }
    if (array_key_exists("HTTP_TAILSCALE_USER_NAME", $_SERVER)) {
        $name = $_SERVER["HTTP_TAILSCALE_USER_NAME"];
    }
    if (array_key_exists("HTTP_TAILSCALE_USER_PROFILE_PIC", $_SERVER)) {
        $profile_pic = $_SERVER["HTTP_TAILSCALE_USER_PROFILE_PIC"];
    }
    if (!defined("TS_LOGIN") || TS_LOGIN == "tagged-devices" || TS_LOGIN == "") {
        $login = "anonymous";
        $name = "Anônimo";
        $profile_pic = "data:image/png;base64,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";
        $has_login = false;
    }
    // define("TS_LOGIN", $login);
    define("TS_NAME", $name);
    define("TS_PROFILE_PIC", $profile_pic);
    define("TS_HAS_LOGIN", $has_login);
}
auth_tailscale();

function rickroll_unlogged() {
    if (!TS_HAS_LOGIN) {
        http_response_code(307);
        set_header("Location", "https://www.youtube.com/watch?v=dQw4w9WgXcQ");
        exit(0); // dispara shutdown, que vai enviar o que precisa ser enviado
    }
}

content_scope_push(); // saporra appenda os echo num buffer pq nessa fase ainda não tem nada retornado

// ==================================== Finalização =========================

function shutdown() {
    global $_HEADERS;
    global $_HEADERS_KV;
    $data = content_scope_pop();

    if (!http_response_code()) {
        http_response_code(200); // default response code
    }
    echo "HTTP/1.0 ";
    echo http_response_code();
    echo "\r\n";

    foreach ($_HEADERS as $header) {
        echo "$header\r\n";
    }

    if ($_HEADERS_KV['Content-Type'] == 'auto') {
        set_contenttype(mime_from_buffer($data));
    }

    foreach ($_HEADERS_KV as $key => $value) {
        echo "$key: $value\r\n";
    }

    echo "\r\n";

    echo $data;

}
register_shutdown_function('shutdown');

// ==================================== ROTAS ===============================

$SCRIPT_DIR = getenv("SCRIPT_DIR");
chdir($SCRIPT_DIR);

$ROUTES_SCRIPT = "$SCRIPT_DIR/routes.php";

if (!file_exists($ROUTES_SCRIPT)) {
    http_response_code(404);
} else {
    include "$SCRIPT_DIR/routes.php";
}

?>
