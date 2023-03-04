import { serve } from "https://deno.land/std@0.178.0/http/server.ts";

const abortController = new AbortController();
Deno.addSignalListener("SIGINT", () => {
  console.log("SIGINT received, shutting down...");
  abortController.abort();
});
Deno.addSignalListener("SIGTERM", () => {
  console.log("SIGTERM received, shutting down...");
  abortController.abort();
});

serve(handler, { signal: abortController.signal });


function handler(req: Request): Response {
  let body = req.method + " " + req.url + "\n";
  for (const [key, value] of req.headers) {
    body += key + ": " + value + "\n";
  }
  return new Response(body);
}
