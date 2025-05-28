import requests
import unittest
import json
import random
import string
import logging

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

BASE_URL = "http://srv17.mikr.us:40331"


def random_string(length=8):
    return ''.join(random.choices(string.ascii_lowercase, k=length))


class CoffeeApiTests(unittest.TestCase):
    def setUp(self):
        self.register_data = {
            "username": random_string(),
            "password": random_string(),
            "email":    f"{random_string()}@example.com"
        }
        reg_resp = requests.post(f"{BASE_URL}/register", json=self.register_data)
        self.assertEqual(reg_resp.status_code, 200, "Registration failed")
        user_json = reg_resp.json()
        self.user_id = user_json.get("id") or user_json.get("user_id")

        login_payload = {
            "username":  self.register_data["username"],
            "passwords": self.register_data["password"]
        }
        login_resp = requests.post(f"{BASE_URL}/login", json=login_payload)
        self.assertEqual(login_resp.status_code, 200, "Login failed")
        token = login_resp.json().get("token")
        self.assertIsNotNone(token, "No token returned")
        self.auth_headers = {"Authorization": f"Bearer {token}"}

        self.coffee_data = {
            "name":         f"Test Coffee {random_string()}",
            "roasteryId":   random.randint(1, 10),
            "country":      random_string(),
            "region":       random_string(),
            "farm":         random_string(),
            "variety":      random_string(),
            "process":      random_string(),
            "roastProfile": random_string(),
            "flavourNotes": [random_string(), random_string(), random_string()],
            "description":  random_string(50),
            "imageUrl":     f"https://example.com/{random_string()}.jpg",
        }
        self.roastery_data = {
            "name":        f"Roastery {random_string()}",
            "country":     random_string(),
            "city":        random_string(),
            "address":     random_string(),
            "website":     f"https://{random_string()}.com/",
            "description": random_string(50),
            "imageUrl":    f"https://example.com/{random_string()}.jpg",
        }
        self.shop_data = {
            "name":        f"Shop {random_string()}",
            "country":     random_string(),
            "city":        random_string(),
            "address":     random_string(),
            "website":     f"https://{random_string()}.com/",
            "description": random_string(50),
            "imageUrl":    f"https://example.com/{random_string()}.jpg",
        }
        self.review_data_coffee = {
            "coffeeId": random.randint(1, 10),
            "rating":   random.randint(1, 5),
            "review":   random_string(30),
        }
        self.review_data_roastery = {
            "roasteryId": random.randint(1, 10),
            "rating":     random.randint(1, 5),
            "review":     random_string(30),
        }
        self.review_data_shop = {
            "coffeeShopId": random.randint(1, 10),
            "rating":       random.randint(1, 5),
            "review":       random_string(30),
        }


    def test_api_documentation(self):
        resp = requests.get(f"{BASE_URL}/")
        self.assertEqual(resp.status_code, 200, "Failed to get API documentation")
        self.assertIn("application/json", resp.headers.get("Content-Type", ""), "Response should be JSON")


    def test_html_documentation(self):
        resp = requests.get(f"{BASE_URL}/help")
        self.assertEqual(resp.status_code, 200, "Failed to get HTML documentation")
        self.assertTrue(resp.text.strip().lower().startswith("<!doctype html>"), "Response should be HTML document")


    def test_register(self):
        payload = {
            "username": f"new_{random_string()}",
            "password": "Password123!",
            "email":    f"{random_string()}@example.com"
        }
        resp = requests.post(f"{BASE_URL}/register", json=payload)
        self.assertEqual(resp.status_code, 200)


    def test_login(self):
        payload = {
            "username":  self.register_data["username"],
            "passwords": self.register_data["password"]
        }
        resp = requests.post(f"{BASE_URL}/login", json=payload)
        self.assertEqual(resp.status_code, 200)
        self.assertIn("token", resp.json())


    def test_get_user_by_id(self):
        resp = requests.get(f"{BASE_URL}/users/{self.user_id}", headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200)

        body = resp.json()
        user = body.get("data") or body.get("user") or body

        self.assertEqual(user.get("username"), self.register_data["username"], "Username should match registered user")


    def test_get_coffees(self):
        resp = requests.get(f"{BASE_URL}/coffees")
        self.assertEqual(resp.status_code, 200, "Failed to get coffees")
        self.assertIsInstance(resp.json(), list, "Response should be a list of coffees")


if __name__ == "__main__":
    unittest.main()