local context = Context.new()

-- Identity
require("author").prompt(context)
require("org").prompt(context)

context:set("suffix_options", { "Service", "Orchestrator", "Adapter", "Router", "Gateway" })
context:set("suffix_default", "Service")
require("project").prompt(context)

context:set("repo_name", context:get("project-name"))
context:set("github_owner", context:get("org-solution-name"))

-- Language-specific extra: Go module path
context:prompt_text("Go Module Path:", "module_path", {
    default = "github.com/" .. context:get("org-solution-name") .. "/" .. context:get("project-name"),
    help = "Go module path (e.g. github.com/" .. context:get("org-solution-name") .. "/billing-service)",
})

-- Service configuration
require("ports").prompt(context, { ports = { { "service", help = "HTTP port for the service" }, "management", "debug" } })

-- EditorConfig + gitignore (pre-seeded; library skips interactive prompt)
local editor_config = require("editor-config")
editor_config.prompt(context, {
    languages     = { "Go", "YAML", "Markdown" },
    gitattributes = true,
})

local gitignore = require("gitignore")
gitignore.prompt(context, {
    ignores = { "Go", "Claude", "IDEA", "VSCode", "macOS" },
})

-- SCM (opt-in: default = None)
local scm = require("scm")
scm.prompt(context)

if archetype.switches.is_enabled("debug-context") then
    log.info(archetype.description .. " Context:")
    output.print(format.yaml(context))
end

-- Render base workspace
directory.render("contents/base", context)

local dest = { destination = context:get("project-name") }

-- CI workflows
local ci = require("golang-ci")
ci.render(context, dest)

-- Platform manifests
-- A basic service weaves in no resources; pre-set these so the manifests library
-- omits resourceRequirements and does not prompt for them.
context:set("persistence", "None")
context:set("cache", "None")
context:set("messaging", "None")
context:set("messaging_access", "produce")
context:set("protocol", "REST")
local platform = require("platform-application-manifests")
platform.prompt(context)
platform.finalize(context, dest)

-- EditorConfig, gitignore, SCM finalize (side effects last)
editor_config.finalize(context, dest)
gitignore.finalize(context, dest)
scm.finalize(context)

-- Archive (zip / tarball switches for Ybor Studio)
require("archiver").finalize(context)

return context
