from pyinfra import host

__ALL__ = []

def expose():
    def handle(fn):
        __ALL__.append(fn.__name__)
        return fn
    return handle

@expose()
def is_ssh():
    return host.executor.__name__ == "pyinfra.connectors.ssh"

@expose()
def is_local():
    return host.executor.__name__ == "pyinfra.connectors.local"

@expose()
def is_nixos():
    from pyinfra.facts.server import LinuxName
    distro = host.get_fact(LinuxName)
    return distro == "NixOS"

