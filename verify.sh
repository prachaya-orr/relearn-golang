#!/bin/bash

BASE_URL="http://localhost:8080"

echo "1. Creating a Todo..."
RESPONSE=$(curl -s -X POST $BASE_URL/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go", "description": "Finish the CRUD tutorial"}')
echo "Response: $RESPONSE"

# Extract ID using python3 (assuming it is available as used below)
# We handle cases where response might be an error or unexpected format loosely for this dev script.
if command -v python3 &> /dev/null; then
    TODO_ID=$(echo $RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin).get('data', {}).get('id', ''))")
else
    # Fallback to simple grep/sed for ID if python3 not found (though less reliable for JSON)
    # This might need to be smarter for nested json, but sticking to simple for now
    TODO_ID=$(echo $RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
fi

if [ -z "$TODO_ID" ]; then
    echo "Error: Could not extract Todo ID from response."
    exit 1
fi

echo "Created Todo ID: $TODO_ID"
echo -e "\n"

echo "2. Listing Todos..."
curl -s $BASE_URL/todos | python3 -m json.tool
echo -e "\n"

echo "3. Updating Todo (ID: $TODO_ID)..."
curl -X PUT $BASE_URL/todos/$TODO_ID \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go (Updated)", "description": "Master the CRUD tutorial", "completed": true}'
echo -e "\n"

echo "4. Getting Todo (ID: $TODO_ID)..."
curl -s $BASE_URL/todos/$TODO_ID | python3 -m json.tool
echo -e "\n"

echo "5. Deleting Todo (ID: $TODO_ID)..."
curl -X DELETE $BASE_URL/todos/$TODO_ID
echo -e "\n"

echo "6. Verifying Deletion..."
curl -s $BASE_URL/todos
echo -e "\nAll Done!"
