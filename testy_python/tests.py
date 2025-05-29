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
        
    
    def test_create_coffee(self):
        resp = requests.post(f"{BASE_URL}/coffees", json=self.coffee_data, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to create coffee")
        coffee = resp.json()
        self.coffee_id_to_update = coffee.get("id") or coffee.get("coffeeId")
        self.assertIn("id", coffee, "Response should contain coffee ID")
        self.assertEqual(coffee["name"], self.coffee_data["name"], "Coffee name should match input data")
        
    
    def test_update_coffee(self):
        if not hasattr(self, 'coffee_id_to_update'):
            self.test_create_coffee()
        
        update_data = {
            "name": f"Updated Coffee {random_string()}",
            "description": random_string(50)
        }
        resp = requests.put(f"{BASE_URL}/coffees/{self.coffee_id_to_update}", json=update_data, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to update coffee")
        updated_coffee = resp.json()
        self.assertEqual(updated_coffee["name"], update_data["name"], "Coffee name should be updated")
        
    
    def test_delete_coffee(self):
        if not hasattr(self, 'coffee_id_to_update'):
            self.test_create_coffee()
        
        resp = requests.delete(f"{BASE_URL}/coffees/{self.coffee_id_to_update}", headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to delete coffee")
        self.assertIn("message", resp.json(), "Response should contain a message")
        
        
    def test_get_roasteries(self):
        resp = requests.get(f"{BASE_URL}/roasteries")
        self.assertEqual(resp.status_code, 200, "Failed to get roasteries")
        self.assertIsInstance(resp.json(), list, "Response should be a list of roasteries")
        
    
    def test_create_roastery(self):
        resp = requests.post(f"{BASE_URL}/roasteries", json=self.roastery_data, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to create roastery")
        roastery = resp.json()
        self.roastery_id_to_update = roastery.get("id") or roastery.get("roasteryId")
        self.assertIn("id", roastery, "Response should contain roastery ID")
        self.assertEqual(roastery["name"], self.roastery_data["name"], "Roastery name should match input data")
        
    
    def test_update_roastery(self):
        if not hasattr(self, 'roastery_id_to_update'):
            self.test_create_roastery()
        
        update_data = {
            "name": f"Updated Roastery {random_string()}",
            "description": random_string(50)
        }
        resp = requests.put(f"{BASE_URL}/roasteries/{self.roastery_id_to_update}", json=update_data, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to update roastery")
        updated_roastery = resp.json()
        self.assertEqual(updated_roastery["name"], update_data["name"], "Roastery name should be updated")
        
    
    def test_delete_roastery(self):
        if not hasattr(self, 'roastery_id_to_update'):
            self.test_create_roastery()
        
        resp = requests.delete(f"{BASE_URL}/roasteries/{self.roastery_id_to_update}", headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to delete roastery")
        self.assertIn("message", resp.json(), "Response should contain a message")
        
        
    def test_get_coffee_shops(self):
        resp = requests.get(f"{BASE_URL}/coffee-shops")
        self.assertEqual(resp.status_code, 200, "Failed to get coffee shops")
        self.assertIsInstance(resp.json(), list, "Response should be a list of coffee shops")
    
    
    def test_create_coffee_shop(self):
        resp = requests.post(f"{BASE_URL}/coffee-shops", json=self.shop_data, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to create coffee shop")
        shop = resp.json()
        self.shop_id_to_update = shop.get("id") or shop.get("coffeeShopId")
        self.assertIn("id", shop, "Response should contain coffee shop ID")
        self.assertEqual(shop["name"], self.shop_data["name"], "Coffee shop name should match input data")
        
        
    def test_update_coffee_shop(self):
        if not hasattr(self, 'shop_id_to_update'):
            self.test_create_coffee_shop()
        
        update_data = {
            "name": f"Updated Coffee Shop {random_string()}",
            "description": random_string(50)
        }
        resp = requests.put(f"{BASE_URL}/coffee-shops/{self.shop_id_to_update}", json=update_data, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to update coffee shop")
        updated_shop = resp.json()
        self.assertEqual(updated_shop["name"], update_data["name"], "Coffee shop name should be updated")
        
        
    def test_delete_coffee_shop(self):
        if not hasattr(self, 'shop_id_to_update'):
            self.test_create_coffee_shop()
        
        resp = requests.delete(f"{BASE_URL}/coffee-shops/{self.shop_id_to_update}", headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to delete coffee shop")
        self.assertIn("message", resp.json(), "Response should contain a message")
        
        
    
    def test_post_coffee_review(self):
        if not hasattr(self, 'coffee_id_to_update'):
            self.test_create_coffee()
        
        self.review_data_coffee["coffeeId"] = self.coffee_id_to_update
        resp = requests.post(f"{BASE_URL}/reviews/coffees", json=self.review_data_coffee, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to post coffee review")
        review = resp.json()
        self.assertIn("id", review, "Response should contain review ID")
        
        
    def test_delete_coffee_review(self):
        if not hasattr(self, 'coffee_id_to_update'):
            self.test_post_coffee_review()
        
        resp = requests.delete(f"{BASE_URL}/reviews/coffees/{self.review_data_coffee['coffeeId']}", headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to delete coffee review")
        self.assertIn("message", resp.json(), "Response should contain a message")
        
        
    
    def test_post_roastery_review(self):
        if not hasattr(self, 'roastery_id_to_update'):
            self.test_create_roastery()
        
        self.review_data_roastery["roasteryId"] = self.roastery_id_to_update
        resp = requests.post(f"{BASE_URL}/reviews/roasteries", json=self.review_data_roastery, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to post roastery review")
        review = resp.json()
        self.assertIn("id", review, "Response should contain review ID")
        
        
    
    def test_delete_roastery_review(self):
        if not hasattr(self, 'roastery_id_to_update'):
            self.test_post_roastery_review()
        
        resp = requests.delete(f"{BASE_URL}/reviews/roasteries/{self.review_data_roastery['roasteryId']}", headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to delete roastery review")
        self.assertIn("message", resp.json(), "Response should contain a message")
        
        
    def test_post_coffee_shop_review(self):
        if not hasattr(self, 'shop_id_to_update'):
            self.test_create_coffee_shop()
        
        self.review_data_shop["coffeeShopId"] = self.shop_id_to_update
        resp = requests.post(f"{BASE_URL}/reviews/coffee-shops", json=self.review_data_shop, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to post coffee shop review")
        review = resp.json()
        self.assertIn("id", review, "Response should contain review ID")
        
        
    def test_delete_coffee_shop_review(self):
        if not hasattr(self, 'shop_id_to_update'):
            self.test_post_coffee_shop_review()
        
        resp = requests.delete(f"{BASE_URL}/reviews/coffee-shops/{self.review_data_shop['coffeeShopId']}", headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, "Failed to delete coffee shop review")
        self.assertIn("message", resp.json(), "Response should contain a message")
    
        
    



if __name__ == "__main__":
    unittest.main()