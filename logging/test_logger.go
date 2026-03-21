package logging



func ExampleLogger() {
	cfg := NewLoggingConfig(true, "DEBUG", "txt", 1)
	bl := NewBaseLogger("0.0.1", "C:\\Users\\itursunov\\Documents\\GitlabProjects\\driverlog\\logs", cfg)

	bl.InvokeLogging()

	logger := NewModuleLogger("TEST", "1.0", bl)
	logger2 := NewModuleLogger("TEST", "2.0", logger)

	logger.Debug("DEBUG MESSAGE")
	logger.Info("INFO MESSAGE")
	logger.Warn("WARN MESSAGE")
	logger.Error("ERROR MESSAGE")
	logger2.Debug("DEBUG MESSAGE")
	logger2.Critical("CRITICAL MESSAGE")
}