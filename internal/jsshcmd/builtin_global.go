package jsshcmd

import (
	"fmt"
	"github.com/leizongmin/go/configloader"
	_ "github.com/leizongmin/go/configloader/toml"
	_ "github.com/leizongmin/go/configloader/yaml"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func jsFnExit(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			os.Exit(0)
			return ctx.Int32(0)
		}
		if !args[0].IsNumber() {
			return ctx.ThrowTypeError("exit: first argument expected number type")
		}
		code := args[0].Int32()
		os.Exit(int(code))
		return ctx.Int32(code)
	}
}

func jsFnSleep(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("sleep: missing millisecond argument")
		}
		if !args[0].IsNumber() {
			return ctx.ThrowTypeError("sleep: first argument expected number type")
		}
		ret, err := jsexecutor.JSValueToAny(args[0])
		if err != nil {
			return ctx.ThrowError(err)
		}
		ms, ok := ret.(float64)
		if !ok {
			return ctx.ThrowTypeError("sleep: first argument expected number type")
		}

		time.Sleep(time.Millisecond * time.Duration(ms))
		return jsexecutor.AnyToJSValue(ctx, ms)
	}
}

func isSupportedConfigFormat(t string) bool {
	if t == "json" || t == "yaml" || t == "toml" {
		return true
	}
	return false
}

func jsFnLoadconfig(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("loadconfig: missing name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("loadconfig: first argument expected string type")
		}
		file := args[0].String()

		format := ""
		if len(args) >= 2 {
			if !args[1].IsString() {
				return ctx.ThrowTypeError("loadconfig: second argument expected string type")
			}
			format = strings.ToLower(args[1].String())
			if !isSupportedConfigFormat(format) {
				return ctx.ThrowTypeError("loadconfig: second argument only accepted one of json,yaml,toml")
			}
		}

		ext := strings.ToLower(filepath.Ext(file))
		if ext == "" {
			ext = ".json"
		} else if ext == ".yml" {
			ext = ".yaml"
		}
		if format == "" {
			format = ext[1:]
		}

		content, err := ioutil.ReadFile(file)
		if err != nil {
			return ctx.ThrowError(err)
		}

		data := make(typeutil.H)
		if err := configloader.Load(format, content, &data); err != nil {
			return ctx.ThrowError(err)
		}

		return jsexecutor.AnyToJSValue(ctx, data)
	}
}

func jsFnReadline(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		var line string
		_, err := fmt.Scanln(&line)
		if err != nil {
			if err.Error() == "unexpected newline" {
				return ctx.String("")
			}
			return ctx.ThrowError(err)
		}
		return ctx.String(line)
	}
}

func jsFnEvalfile(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("evalfile: missing filename")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("evalfile: first argument expected string type")
		}
		file := args[0].String()

		var content string
		if len(args) >= 2 {
			if !args[1].IsString() {
				return ctx.ThrowTypeError("evalfile: seconds argument expected string type")
			}
			content = args[1].String()
		} else if isUrl(file) {
			s, err := httpGetFileContent(file)
			if err != nil {
				return ctx.ThrowError(err)
			}
			content = s
		} else {
			b, err := ioutil.ReadFile(file)
			if err != nil {
				return ctx.ThrowError(err)
			}
			content = string(b)
		}

		ret, err := ctx.EvalFile(content, file)
		if err != nil {
			return ctx.ThrowError(err)
		}
		return ret
	}
}

func jsFnBytesize(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("bytesize: missing data")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("bytesize: first argument expected string type")
		}
		data := args[0].String()

		b := []byte(data)
		return ctx.Int32(int32(len(b)))
	}
}
