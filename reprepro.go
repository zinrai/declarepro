package main

func (r *Runner) repreproIncludeDeb(baseDir, dist, file string) error {
	_, err := r.run("reprepro", "-b", baseDir, "includedeb", dist, file)
	return err
}
