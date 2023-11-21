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

def test_suite1():
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
    id = None
    try:
        client.register(test_user)
        client.signin(test_user)
        id = client.add_item(test_item)
        print(id)
        item = client.get_item(id)
        print(item)
    finally:
        client.remove_item(id)
        client.delete_account(test_user)

test_suite1()
