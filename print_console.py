from fastapi import FastAPI, Request, UploadFile
import json

app = FastAPI()

@app.post("/print")
async def print_any(request: Request):
    headers = dict(request.headers)
    query_params = dict(request.query_params)
    body_repr = None
    body_type = None

    try:
        body_json = await request.json()
        body_repr = body_json
        body_type = "json"
    except:
        body_bytes = await request.body()
        try:
            body_repr = body_bytes.decode("utf-8")
            body_type = "text"
        except:
            body_repr = body_bytes
            body_type = "bytes"

    print("=== /print received ===")
    print("Headers:", json.dumps(headers, indent=2, ensure_ascii=False))
    print("Query params:", json.dumps(query_params, indent=2, ensure_ascii=False))
    print("Body type:", body_type)
    print("Body:", body_repr)
    print("=======================")

    return {"ok": True, "body_type": body_type}
