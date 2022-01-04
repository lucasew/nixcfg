local fennel = require 'fennel'

print(#package.loaders)
table.insert(_G.package.loaders or _G.package.searchers, fennel.searcher)
print(#package.loaders)

fennel.path = package.path:gsub(".lua", ".fnl")

function getScriptFile()
    return debug.getinfo(2, 'S').source:sub(2)
end

print(getScriptFile())

return require 'adapter'

