package running

import (
    r "git.masontest.com/branches/goserver/core/running"
    "testing"
    "git.masontest.com/branches/goserver/plugins"
    "git.masontest.com/branches/goserver/config"
)

func Test_Build(t *testing.T) {
    logger := plugins.NewServerLog("test")

    _, serviceConfig   := config.GetConfig("../../config.yml", "dev")

    for appName, config := range serviceConfig {
        r.Build(appName, config["handle_path"], logger)
    }
}
