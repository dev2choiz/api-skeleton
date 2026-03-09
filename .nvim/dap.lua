local configurations = {
  go = {
    {
      -- Override this configuration to avoid being asked for the path
      name = "Attach Go (in docker)",
      type = "go",
      request = "attach",
      mode = "remote",
      substitutePath = {
        {
          from = "${workspaceFolder}",
          to = "/app",
        },
      },
    },
  },
}

return {
  adapters = {
    -- Override this adapter to avoid being asked for the delve port
    go = {
      type = "server",
      host = "localhost",
      port = "2345",
    },
  },
  configurations = configurations,
}
