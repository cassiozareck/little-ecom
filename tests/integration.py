import requests

# Constants for the API endpoints
BASE_URL = 'http://192.168.49.2'
REGISTER_URL = f'{BASE_URL}/auth/register'
SIGNIN_URL = f'{BASE_URL}/auth/signin'
ADD_ITEM_URL = f'{BASE_URL}/item'

# Test data
test_user = {
    "email": "testuser@example.com",
    "password": "securepassword"
}

test_item = {
    "name": "Test Item",
    "owner": "testOwnerId",  # Replace with an actual owner ID after registration
    "price": 9.99
}

def test_register():
    response = requests.post(REGISTER_URL, json=test_user)
    assert response.status_code == 201

def test_signin():
    response = requests.post(SIGNIN_URL, json=test_user)
    assert response.status_code == 200
    token = response.text
    return token

def test_add_item():
    token = test_signin()  # Get a valid token
    headers = {'Authorization': f'Bearer {token}'}
    response = requests.post(ADD_ITEM_URL, json=test_item, headers=headers)
    assert response.status_code == 200

# Additional test functions for other endpoints

def test_delete_item():
    token = test_signin()  # Get a valid token
    item_id = "testItemId"  # Replace with an actual item ID after adding an item
    headers = {'Authorization': f'Bearer {token}'}
    response = requests.delete(f'{BASE_URL}/item/{item_id}', headers=headers)
    assert response.status_code == 200

def test_update_item():
    token = test_signin()  # Get a valid token
    item_id = "testItemId"  # Replace with an actual item ID after adding an item
    headers = {'Authorization': f'Bearer {token}'}
    updated_item = {
        "name": "Updated Test Item",
        "owner": "testOwnerId",  # Replace with an actual owner ID
        "price": 19.99
    }
    response = requests.put(f'{BASE_URL}/item/{item_id}', json=updated_item, headers=headers)
    assert response.status_code == 200

def test_get_item_by_id():
    item_id = "testItemId"  # Replace with an actual item ID after adding an item
    response = requests.get(f'{BASE_URL}/item/{item_id}')
    assert response.status_code == 200

def test_get_items_by_owner():
    token = test_signin()  # Get a valid token
    owner_id = "testOwnerId"  # Replace with an actual owner ID
    headers = {'Authorization': f'Bearer {token}'}
    response = requests.get(f'{BASE_URL}/items/{owner_id}', headers=headers)
    assert response.status_code == 200

def test_get_all_items():
    response = requests.get(f'{BASE_URL}/items')
    assert response.status_code == 200

# Test execution order
def test_api_workflow():
    test_register()
    test_signin()
    test_add_item()
    test_get_item_by_id()
    test_get_items_by_owner()
    test_get_all_items()
    test_update_item()
    test_delete_item()

test_api_workflow()