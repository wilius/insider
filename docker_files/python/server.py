from flask import Flask, request, jsonify
import re
import secrets
import string

app = Flask(__name__)

# Regular expression to match 4xx and 5xx status codes
STATUS_CODE_PATTERN = re.compile(r'^(4\d\d|5\d\d)$')

def generate_random_string(length=16):
    characters = string.ascii_letters + string.digits
    return ''.join(secrets.choice(characters) for _ in range(length))

@app.route('/webhook', methods=['POST'])
def webhook_handler():
    try:
        data = request.json
        if not data or 'content' not in data:
            return jsonify({"error": "Invalid payload"}), 400

        # Extract content and check if it's a status code
        content = str(data.get('content'))
        if STATUS_CODE_PATTERN.match(content):
            # Return the extracted status code as an integer
            return jsonify({"message": f"Responding with status {content}"}), int(content)

        # If content is not a 4xx or 5xx code, return 200
        return jsonify({"message": "Success", "messageId": generate_random_string() }), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
