import requests

class APITestClient:
    BASE_URL = 'http://192.168.49.2'

    def __init__(self):
        self.token = None

    def register(self, user_data):
        response = requests.post(f'{self.BASE_URL}/auth/register', json=user_data)
        print("register: ", response, response.text)
        assert response.status_code == 201

    def signin(self, user_data):
        response = requests.post(f'{self.BASE_URL}/auth/signin', json=user_data)
        print("signin: ", response, response.text)
        assert response.status_code == 200
        self.token = response.text

    def add_item(self, item_data):
        headers = {'Authorization': f'Bearer {self.token}'}
        response = requests.post(f'{self.BASE_URL}/item', json=item_data, headers=headers)
        print("add_item: ", response, response.text)
        assert response.status_code == 201
        return response.text

    def get_item(self, item_id):
        headers = {'Authorization': f'Bearer {self.token}'}
        response = requests.get(f'{self.BASE_URL}/item/{item_id}', headers=headers)
        print("get_item: ", response, response.text)
        assert response.status_code == 200
        return response.json()

    def update_item(self, item_id, item_data):
        headers = {'Authorization': f'Bearer {self.token}'}
        response = requests.put(f'{self.BASE_URL}/item/{item_id}', json=item_data, headers=headers)
        print("update_item: ", response, response.text)
        assert response.status_code == 200

    def get_items(self):
        headers = {'Authorization': f'Bearer {self.token}'}
        response = requests.get(f'{self.BASE_URL}/items', headers=headers)
        print("get_items: ", response, response.text)
        assert response.status_code == 200
        return response.json()

    def get_items_by_owner(self, owner):
        headers = {'Authorization': f'Bearer {self.token}'}
        response = requests.get(f'{self.BASE_URL}/items/{owner}', headers=headers)
        print("get_items_by_owner: ", response, response.text)
        assert response.status_code == 200
        return response.json()

    # should pass the account to be deleted in body (email and password)
    def delete_account(self, account):
        response = requests.delete(f'{self.BASE_URL}/auth/delete', json=account)
        print("delete_account: ", response, response.text)
        assert response.status_code == 204

    def remove_item(self, item_id):
        headers = {'Authorization': f'Bearer {self.token}'}
        response = requests.delete(f'{self.BASE_URL}/item/{item_id}', headers=headers)
        print("remove_item: ", response, response.text)
        assert response.status_code == 200

# Basic test suite that will register a user, add an item, update the item, and
# delete the item and user, not making any mistake.
def test_suite1():
    print()
    print("TEST SUITE 1: --------")
    print()

    client = APITestClient()
    test_user = {
        "email": "testowner@example.com",
        "password": "securepassword"
    }
    test_item = {
        "name": "Test Item",
        "owner": "testowner",
        "price": 9.99
    }

    try:
        client.register(test_user)
        client.signin(test_user)
        id = client.add_item(test_item)
        item = client.get_item(id)
        item["price"] = 19.99
        client.update_item(id, item)
        item = client.get_item(id)
        assert item["name"] == "Test Item"
        assert item["owner"] == "testowner"
        assert item["price"] == 19.99

    finally:
        client.remove_item(id)
        client.delete_account(test_user)


# This test suite will add two items and check if get_items endpoints return
# the correct items.
def test_suite2():
    print()
    print("TEST SUITE 2: --------")
    print()

    client = APITestClient()
    test_user = {
        "email": "megy@gmail.com",
        "password": "securepassword"
    }
    test_item = {
        "name": "apple",
        "owner": "megy",
        "price": 0.99
    }
    test_item2 = {
        "name": "pineapple",
        "owner": "megy",
        "price": 1.99
    }

    try:
        client.register(test_user)
        client.signin(test_user)
        id = client.add_item(test_item)
        id2 = client.add_item(test_item2)
        items = client.get_items()
        items2 = client.get_items_by_owner("megy")

        # Check if all items in items2 are in items
        assert all(item in items for item in items2), "Not all items from items2 are in items"

    finally:
        client.remove_item(id)
        client.remove_item(id2)
        client.delete_account(test_user)

# This test suite will try to register a user with an invalid email and password
def test_suite3():
    print()
    print("TEST SUITE 3: --------")
    print()

    client = APITestClient()
    test_user = {
        "email": "",
        "password": ""
    }
    try:
        client.register(test_user)
    except AssertionError:
        print("Response expected since email and password are invalid")

# Test suite 4 will try to add an item without signing in
def test_suite4():
    print()
    print("TEST SUITE 4: --------")
    print()

    client = APITestClient()
    test_item = {
        "name": "Test Item",
        "owner": "testowner",
        "price": 9.99
    }
    try:
        client.add_item(test_item)
    except AssertionError:
        print("Response expected since user is not signed in")


if __name__ == '__main__':
    test_suite1()
    test_suite2()
    test_suite3()
    test_suite4()
