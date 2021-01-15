(()=>{if("undefined"!=typeof global);else if("undefined"!=typeof window)window.global=window;else{if("undefined"==typeof self)throw new Error("cannot export Go (neither global, window nor self is defined)");self.global=self}global.require||"undefined"==typeof require||(global.require=require),!global.fs&&global.require&&(global.fs=require("fs"));const e=()=>{const e=new Error("not implemented");return e.code="ENOSYS",e};if(!global.fs){let t="";global.fs={constants:{O_WRONLY:-1,O_RDWR:-1,O_CREAT:-1,O_TRUNC:-1,O_APPEND:-1,O_EXCL:-1},writeSync(e,o){const n=(t+=s.decode(o)).lastIndexOf("\n");return-1!=n&&(console.log(t.substr(0,n)),t=t.substr(n+1)),o.length},write(t,s,o,n,i,r){if(0!==o||n!==s.length||null!==i)return void r(e());r(null,this.writeSync(t,s))},chmod(t,s,o){o(e())},chown(t,s,o,n){n(e())},close(t,s){s(e())},fchmod(t,s,o){o(e())},fchown(t,s,o,n){n(e())},fstat(t,s){s(e())},fsync(e,t){t(null)},ftruncate(t,s,o){o(e())},lchown(t,s,o,n){n(e())},link(t,s,o){o(e())},lstat(t,s){s(e())},mkdir(t,s,o){o(e())},open(t,s,o,n){n(e())},read(t,s,o,n,i,r){r(e())},readdir(t,s){s(e())},readlink(t,s){s(e())},rename(t,s,o){o(e())},rmdir(t,s){s(e())},stat(t,s){s(e())},symlink(t,s,o){o(e())},truncate(t,s,o){o(e())},unlink(t,s){s(e())},utimes(t,s,o,n){n(e())}}}if(global.process||(global.process={getuid:()=>-1,getgid:()=>-1,geteuid:()=>-1,getegid:()=>-1,getgroups(){throw e()},pid:-1,ppid:-1,umask(){throw e()},cwd(){throw e()},chdir(){throw e()}}),!global.crypto){const e=require("crypto");global.crypto={getRandomValues(t){e.randomFillSync(t)}}}global.performance||(global.performance={now(){const[e,t]=process.hrtime();return 1e3*e+t/1e6}}),global.TextEncoder||(global.TextEncoder=require("util").TextEncoder),global.TextDecoder||(global.TextDecoder=require("util").TextDecoder);const t=new TextEncoder("utf-8"),s=new TextDecoder("utf-8");var o=[];if(global.Go=class{constructor(){this._callbackTimeouts=new Map,this._nextCallbackTimeoutID=1;const e=()=>new DataView(this._inst.exports.memory.buffer),n=(t,s)=>{e().setUint32(t+0,s,!0),e().setUint32(t+4,Math.floor(s/4294967296),!0)},i=t=>{const s=e().getFloat64(t,!0);if(0===s)return;if(!isNaN(s))return s;const o=e().getUint32(t,!0);return this._values[o]},r=(t,s)=>{if("number"==typeof s)return isNaN(s)?(e().setUint32(t+4,2146959360,!0),void e().setUint32(t,0,!0)):0===s?(e().setUint32(t+4,2146959360,!0),void e().setUint32(t,1,!0)):void e().setFloat64(t,s,!0);switch(s){case void 0:return void e().setFloat64(t,0,!0);case null:return e().setUint32(t+4,2146959360,!0),void e().setUint32(t,2,!0);case!0:return e().setUint32(t+4,2146959360,!0),void e().setUint32(t,3,!0);case!1:return e().setUint32(t+4,2146959360,!0),void e().setUint32(t,4,!0)}let o=this._ids.get(s);void 0===o&&(void 0===(o=this._idPool.pop())&&(o=this._values.length),this._values[o]=s,this._goRefCounts[o]=0,this._ids.set(s,o)),this._goRefCounts[o]++;let n=1;switch(typeof s){case"string":n=2;break;case"symbol":n=3;break;case"function":n=4}e().setUint32(t+4,2146959360|n,!0),e().setUint32(t,o,!0)},l=(e,t,s)=>new Uint8Array(this._inst.exports.memory.buffer,e,t),a=(e,t,s)=>{const o=new Array(t);for(let s=0;s<t;s++)o[s]=i(e+8*s);return o},c=(e,t)=>s.decode(new DataView(this._inst.exports.memory.buffer,e,t)),u=Date.now()-performance.now();this.importObject={wasi_unstable:{fd_write:function(t,n,i,r){if(1==t)for(let t=0;t<i;t++){let i=n+8*t,r=e().getUint32(i+0,!0),l=e().getUint32(i+4,!0);for(let t=0;t<l;t++){let n=e().getUint8(r+t);if(13==n);else if(10==n){let e=s.decode(new Uint8Array(o));o=[],console.log(e)}else o.push(n)}}else console.error("invalid file descriptor:",t);return e().setUint32(r,0,!0),0}},env:{"runtime.ticks":()=>u+performance.now(),"runtime.sleepTicks":e=>{setTimeout(this._inst.exports.go_scheduler,e)},"syscall.Exit":e=>{if(!global.process)throw"trying to exit with code "+e;process.exit(e)},"syscall/js.finalizeRef":e=>{console.error("syscall/js.finalizeRef not implemented")},"syscall/js.finalizeRef":t=>{const s=e().getUint32(t+8,!0);if(this._goRefCounts[s]--,0===this._goRefCounts[s]){const e=this._values[s];this._values[s]=null,this._ids.delete(e),this._idPool.push(s)}},"syscall/js.stringVal":(e,t,s)=>{const o=c(t,s);r(e,o)},"syscall/js.valueGet":(e,t,s,o)=>{let n=c(s,o),l=i(t),a=Reflect.get(l,n);r(e,a)},"syscall/js.valueSet":(e,t,s,o)=>{const n=i(e),r=c(t,s),l=i(o);Reflect.set(n,r,l)},"syscall/js.valueDelete":(e,t,s)=>{const o=i(e),n=c(t,s);Reflect.deleteProperty(o,n)},"syscall/js.valueIndex":(e,t,s)=>{r(e,Reflect.get(i(t),s))},"syscall/js.valueSetIndex":(e,t,s)=>{Reflect.set(i(e),t,i(s))},"syscall/js.valueCall":(t,s,o,n,l,u,d)=>{const f=i(s),g=c(o,n),h=a(l,u);try{const s=Reflect.get(f,g);r(t,Reflect.apply(s,f,h)),e().setUint8(t+8,1)}catch(s){r(t,s),e().setUint8(t+8,0)}},"syscall/js.valueInvoke":(t,s,o,n,l)=>{try{const l=i(s),c=a(o,n);r(t,Reflect.apply(l,void 0,c)),e().setUint8(t+8,1)}catch(s){r(t,s),e().setUint8(t+8,0)}},"syscall/js.valueNew":(t,s,o,n,l)=>{const c=i(s),u=a(o,n);try{r(t,Reflect.construct(c,u)),e().setUint8(t+8,1)}catch(s){r(t,s),e().setUint8(t+8,0)}},"syscall/js.valueLength":e=>i(e).length,"syscall/js.valuePrepareString":(e,s)=>{const o=String(i(s)),l=t.encode(o);r(e,l),n(e+8,l.length)},"syscall/js.valueLoadString":(e,t,s,o)=>{const n=i(e);l(t,s).set(n)},"syscall/js.valueInstanceOf":(e,t)=>i(v_attr)instanceof i(t),"syscall/js.copyBytesToGo":(t,s,o,r,a)=>{let c=t,u=t+4;const d=l(s,o),f=i(a);if(!(f instanceof Uint8Array))return void e().setUint8(u,0);const g=f.subarray(0,d.length);d.set(g),n(c,g.length),e().setUint8(u,1)},"syscall/js.copyBytesToJS":(t,s,o,r,a)=>{let c=t,u=t+4;const d=i(s),f=l(o,r);if(!(d instanceof Uint8Array))return void e().setUint8(u,0);const g=f.subarray(0,d.length);d.set(g),n(c,g.length),e().setUint8(u,1)}}}}async run(e){this._inst=e,this._values=[NaN,0,null,!0,!1,global,this],this._goRefCounts=[],this._ids=new Map,this._idPool=[],this.exited=!1;new DataView(this._inst.exports.memory.buffer);for(;;){const e=new Promise(e=>{this._resolveCallbackPromise=(()=>{if(this.exited)throw new Error("bad callback: Go program has already exited");setTimeout(e,0)})});if(this._inst.exports._start(),this.exited)break;await e}}_resume(){if(this.exited)throw new Error("Go program has already exited");this._inst.exports.resume(),this.exited&&this._resolveExitPromise()}_makeFuncWrapper(e){const t=this;return function(){const s={id:e,this:this,args:arguments};return t._pendingEvent=s,t._resume(),s.result}}},global.require&&global.require.main===module&&global.process&&global.process.versions&&!global.process.versions.electron){3!=process.argv.length&&(console.error("usage: go_js_wasm_exec [wasm binary] [arguments]"),process.exit(1));const e=new Go;WebAssembly.instantiate(fs.readFileSync(process.argv[2]),e.importObject).then(t=>e.run(t.instance)).catch(e=>{console.error(e),process.exit(1)})}})();
