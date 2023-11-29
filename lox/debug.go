package lox

import (
	"github.com/rs/zerolog/log"
)

func debugEnterEnv(ctx *Context, name string) func() {
	popEnv := ctx.PushEnv(name)
	log.Debug().Msgf("(%s) ENTER %s", ctx.Phase(), ctx.env)
	return func() {
		log.Debug().Msgf("(%s) EXIT %s", ctx.Phase(), ctx.env)
		popEnv()
	}
}

func debugSetValue(phase Phase, env *Env, name string, to Value) error {
	prev := env.SetValue(name, to)
	if prev == nil {
		log.Debug().Msgf("(%s) SET Env(%s) %s := %s", phase, env.Name(), name, to)
	} else {
		log.Debug().Msgf("(%s) SET Env(%s) %s = %s (was %s)", phase, env.Name(), name, to, prev)
	}
	return nil
}

func debugSetType(phase Phase, env *Env, name string, to Type) error {
	prev := env.SetType(name, to)
	if prev == TypeNone {
		log.Debug().Msgf("(%s) SET Env(%s) %s := %s", phase, env.Name(), name, to)
	} else {
		log.Debug().Msgf("(%s) SET Env(%s) %s = %s (was %s)", phase, env.Name(), name, to, prev)
	}
	return nil
}
