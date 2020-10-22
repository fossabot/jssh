#!/usr/bin/env go run github.com/leizongmin/jssh

log.info("aaa: %s %s %s", __bin, __dirname, __filename)
log.error(JSON.stringify(__args))
println(JSON.stringify(__env))
println("%f", Date.now())

sleep(500)
sh.setenv("__xx__", new Date().toString())
log.info(JSON.stringify(__env))

log.info(sh.pwd())
log.info(sh.cd(__dirname))
log.info(sh.cwd())

sh.exec("pwd")
log.info("%f %f %s", __code, __outputbytes, __output)

sh.exec("pwd", {}, true)
log.info("%f %f %s", __code, __outputbytes, __output)

if (sh.exec(`ls -al ${__homedir}`, {}, true).code === 0) {
    __output.split("\n").forEach(line => log.error(line))
}

fs.readdir(__homedir).forEach(f => log.error(JSON.stringify(f)))
log.info(fs.readfile(`${__homedir}/.gitconfig`))
log.info(JSON.stringify(fs.stat(`${__homedir}/.gitconfig`)))

set("xyz", 12345)
log.info("xyz = %f", xyz)

const file = `${__tmpdir}/${Date.now()}-${Math.random()}.txt`
log.info(fs.writefile(file, "hello"))
log.info(fs.appendfile(file, "world"))
log.info(fs.readfile(file))

log.info(path.abs("."))
log.info(path.base(file))
log.info(path.dir(file))
log.info(path.ext(file))
log.info(path.join("a", "b", "c"))
log.info(path.abs(path.join("a", "b", "c")))

log.info("%s, %s, %v, %v, %s, %s", cli.get(0), cli.get("n"), cli.bool("n"), cli.bool("x"), JSON.stringify(cli.args()), JSON.stringify(cli.opts()))

log.error(new Error().stack)

log.info(JSON.stringify(http.request("GET", "http://baidu.com")))
log.info(format("%s-%s", "aaa", "bbb"))

log.info("bgexec: pid=%v", sh.bgexec("ping qq.com -c 60"))
log.info("bgexec: pid=%v", sh.bgexec("ping baidu.com -c 60"))
log.info("tail: %s", JSON.stringify(sh.exec(`tail ${__filename}`, {}, true)))
sleep(10000)

exit(123)