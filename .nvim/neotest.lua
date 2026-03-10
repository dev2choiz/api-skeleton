return {
  adapters = {
    -- Override this configuration to add `coverprofile` argument
    ["neotest-golang"] = {
      go_test_args = {
        "-v",
        "-race",
        "-count=1",
        "-timeout=60s",
        "-coverprofile=" .. vim.fn.getcwd() .. "/coverage.out",
      },
      dap_go_enabled = true,
      testify_enabled = true,
    },
  },
}
