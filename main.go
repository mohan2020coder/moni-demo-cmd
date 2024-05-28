package main

import (
	"flag"
	"fmt"

	"os"
	"path/filepath"
)

func main() {
	projectName := flag.String("name", "deno-react-app", "Name of the project")
	flag.Parse()

	fmt.Printf("Creating project %s...\n", *projectName)

	createProjectStructure(*projectName)
	createFiles(*projectName)

	fmt.Println("Project created successfully!")
	fmt.Println("Run the following commands to get started:")
	fmt.Printf("cd %s\n", *projectName)
	fmt.Println("deno task dev")
}

func createProjectStructure(projectName string) {
	directories := []string{
		filepath.Join(projectName, "src"),
		projectName,
	}

	fmt.Println(directories)

	for _, dir := range directories {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}
}

func createFiles(projectName string) {
	files := map[string]string{
		filepath.Join(projectName, "deno.json"):       denoConfigContent,
		filepath.Join(projectName, "import_map.json"): importMapContent,
		filepath.Join(projectName, "dev.ts"):          devTsContent,
		filepath.Join(projectName, "index.html"):      indexHtmlContent,
		filepath.Join(projectName, "src", "main.tsx"): mainTsxContent,
		filepath.Join(projectName, "src", "App.tsx"):  appTsxContent,
	}

	for path, content := range files {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", path, err)
			os.Exit(1)
		}
	}
}

const denoConfigContent = `{
  "tasks": {
    "dev": "deno run -A --watch=static/,src/ --unstable dev.ts"
  },
  "importMap": "import_map.json"
}`

const importMapContent = `{
  "imports": {
    "react": "https://esm.sh/react@18.0.0",
    "react-dom": "https://esm.sh/react-dom@18.0.0"
  }
}`

const devTsContent = `import { serve } from "https://deno.land/std@0.140.0/http/server.ts";
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
`

const indexHtmlContent = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Deno React App</title>
</head>
<body>
  <div id="root"></div>
  <script type="module" src="/static/main.js"></script>
</body>
</html>`

const mainTsxContent = `import React from "react";
import ReactDOM from "react-dom";
import App from "./App.tsx";

ReactDOM.render(<App />, document.getElementById("root"));`

const appTsxContent = `import React from "react";

function App() {
  return (
    <div>
      <h1>Hello, Deno with React!</h1>
    </div>
  );
}

export default App;`
