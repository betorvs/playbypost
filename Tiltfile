local_resource(
  name='playbypost-server',
  cmd='task tidy build_assets build_local',
  serve_cmd='./playbypost',
  deps=["app/", "go.mod", "go.sum", "core/"]
)
local_resource(
  name='slack-plugin',
  serve_cmd='./slack-plugin',
  resource_deps=['playbypost-server'],
  deps=["./slack-plugin"]
)