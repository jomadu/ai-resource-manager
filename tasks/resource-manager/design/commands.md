arm version
arm help
arm install
arm outdated [--output <table|json|list>]
arm update
arm uninstall
arm list
arm info
arm clean cache [--nuke]
arm clean sinks

arm add registry --type <git|gitlab|cloudsmith> [--gitlab-group-id ID] [--gitlab-project-id ID] NAME URL
arm remove registry NAME
arm config registry set NAME KEY VALUE
arm list registry
arm info registry [NAME]...

arm add sink [--type <cursor|copilot|amazonq>] [--layout <hierarchical|flat>] [--compile-to <md|cursor|amazonq|copilot>] NAME PATH
arm remove sink NAME
arm config sink set NAME KEY VALUE
arm list sink
arm info sink [NAME]...

arm install ruleset [--priority PRIORITY] [--include GLOB...] [--exclude GLOB...] REGISTRY_NAME/RULESET_NAME[@VERSION] SINK_NAME...
arm uninstall ruleset REGISTRY_NAME/RULESET_NAME
arm config ruleset set REGISTRY_NAME/RULESET_NAME KEY VALUE
arm list ruleset
arm info ruleset [REGISTRY_NAME/RULESET_NAME...]
arm update ruleset [REGISTRY_NAME/RULESET_NAME...]
arm outdated [--output <table|json|list>] ruleset

arm install promptset [--include GLOB...] [--exclude GLOB...] REGISTRY_NAME/PROMPTSET[@VERSION] SINK_NAME...
arm uninstall promptset REGISTRY_NAME/PROMPTSET
arm config promptset set REGISTRY_NAME/PROMPTSET KEY VALUE
arm list promptset
arm info promptset [REGISTRY_NAME/PROMPTSET...]
arm update promptset [REGISTRY_NAME/PROMPTSET...]
arm outdated [--output <table|json|list>] promptset

arm compile [--target <md|cursor|amazonq|copilot>] [--force] [--recursive] [--verbose] [--validate-only] [--include GLOB...] [--output GLOB...] [--fail-fast] INPUT_PATH... OUTPUT_PATH
