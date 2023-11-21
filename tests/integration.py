import requests

class APITestClient:
    BASE_URL = 'http://192.168.49.2'

    def __init__(self):
        self.token = None

    def register(self, user_data):
        response = requests.post(f'{self.BASE_URL}/auth/register', json=user_data)
        assert response.status_code == 201

    def signin(self, user_data):
        response = requests.post(f'{self.BASE_URL}/auth/signin', json=user_data)
        assert response.status_code == 200
        self.token = response.text

    def add_item(self, item_data):
        headers = {'Authorization': f'Bearer {self.token}'}
        response = requests.post(f'{self.BASE_URL}/item', json=item_data, headers=headers)
        assert response.status_code == 200

def test_api_workflow():
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

    client.register(test_user)
    client.signin(test_user)
    client.add_item(test_item)

test_api_workflow()
