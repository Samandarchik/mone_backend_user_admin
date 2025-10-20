from flask import Flask, request

app = Flask(__name__)

@app.route('/agent', methods=['POST'])
def agent():
    try:
        # So‘rov JSON formatidami – tekshiramiz
        data = request.get_json(force=True)
        print("📦 JSON Body:", data)
        return {"status": "success", "received": data}, 200
    except Exception as e:
        # Agar JSON emas, raw bodyni olaylik
        raw_data = request.data.decode('utf-8')
        print("🧾 RAW Body:", raw_data)
        return {"success": True, "raw_received": raw_data}, 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=1212, debug=True)
