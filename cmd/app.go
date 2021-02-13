package cmd

func Run() {
	RootCmd.AddCommand(WalkCommand)
	RootCmd.AddCommand(ExportCmd)
	RootCmd.AddCommand(ExportGearCmd)
	RootCmd.AddCommand(ExportMechCmd)
	RootCmd.AddCommand(ParseCmd)
	RootCmd.AddCommand(ImportCmd)
	RootCmd.AddCommand(LintCmd)
	RootCmd.Execute()
}
