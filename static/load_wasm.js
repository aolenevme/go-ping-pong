!async function(){const a=new Go,t=await fetch("client.wasm"),n=await t.arrayBuffer(),e=await WebAssembly.instantiate(n,a.importObject);a.run(e.instance)}();
