# receives: list of modules to import, type of item looking for in modules
# returns: paths of the modules it found
items:
component:
let
  pathListIfExist = import <dotfiles/lib/pathListIfExist.nix>;
  normalizedComponent = "/" + component;
  suffixifyItem = item: item + normalizedComponent + ".nix";
  suffixedItems = map suffixifyItem items;
  filterExistentItems = items: if (builtins.length items == 0) 
    then [] 
    else 
      pathListIfExist (builtins.head items) ++ (filterExistentItems (builtins.tail items));
  filteredItems = filterExistentItems suffixedItems;
in filteredItems
