local_resource(
  name='playbypost-server',
  cmd='task tidy build_assets build_local',
  serve_cmd='./playbypost -autoplay-worker',
  deps=["app/", "go.mod", "go.sum", "core/"]
)
local_resource(
  name='discord-plugin',
  serve_cmd='./discord-plugin',
  resource_deps=['playbypost-server'],
  deps=["./discord-plugin"]
)