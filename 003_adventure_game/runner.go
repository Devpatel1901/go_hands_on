package main

type Runner interface {
	Start(provider *StoryArcProvider)
}
