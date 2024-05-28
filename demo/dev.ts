import { serve } from "https://deno.land/std@0.140.0/http/server.ts";
import { compile } from "https://deno.land/x/tsc@v0.9.3/mod.ts";

const PORT = 3000;

async function handler(req: Request): Promise<Response> {
  const url = new URL(req.url);
  let path = "." + url.pathname;
  if (path == "./") path = "./index.html";

  try {
    const file = await Deno.readFile(path);
    const contentType = path.endsWith(".html") ? "text/html" :
                        path.endsWith(".js") ? "application/javascript" :
                        path.endsWith(".css") ? "text/css" :
                        "application/octet-stream";
    return new Response(file, { headers: { "Content-Type": contentType } });
  } catch {
    return new Response("404 Not Found", { status: 404 });
  }
}

async function startServer() {
  console.log("Server running on http://localhost:${PORT}/");
  await serve(handler, { addr: ":" + PORT });
}

async function compileAndWatch() {
  await compile({
    entryPoints: ["./src/main.tsx"],
    outDir: "./static",
    compilerOptions: {
      jsx: "react",
      jsxFactory: "React.createElement",
      jsxFragmentFactory: "React.Fragment",
      target: "es2015",
      module: "esnext"
    }
  });

  const watcher = Deno.watchFs(["./src"]);
  for await (const _ of watcher) {
    console.log("Recompiling...");
    await compile({
      entryPoints: ["./src/main.tsx"],
      outDir: "./static",
      compilerOptions: {
        jsx: "react",
        jsxFactory: "React.createElement",
        jsxFragmentFactory: "React.Fragment",
        target: "es2015",
        module: "esnext"
      }
    });
  }
}

await compileAndWatch();
await startServer();
