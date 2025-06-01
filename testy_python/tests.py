import requests
import unittest
import json
import random
import string
import logging
import time

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

BASE_URL = "http://srv17.mikr.us:40331"

def random_string(length=8):
    return ''.join(random.choices(string.ascii_lowercase, k=length))

class CoffeeApiTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.admin_user_id = 1 
        cls.admin_username = "admin"

        admin_credentials = {
            "username": cls.admin_username,
            "passwords": "admin123" 
        }
        
        login_resp = requests.post(f"{BASE_URL}/login", json=admin_credentials)
        if login_resp.status_code != 200:
            logging.error(f"Admin login failed: {login_resp.status_code} - {login_resp.text}")
            logging.error(f"Admin login payload: {admin_credentials}")
        assert login_resp.status_code == 200, f"Admin login failed: {login_resp.text}"
        
        token_data = login_resp.json()
        token = token_data.get("token") or token_data.get("access_token")
        assert token is not None, f"No token returned from admin login: {token_data}"
        cls.auth_headers = {"Authorization": f"Bearer {token}"}
        
        cls.user_id_for_tests = str(cls.admin_user_id)

        cls.new_user_register_data_template = {
            "username": f"testuser_{random_string(10)}",
            "password": "TestPassword123!", 
            "email": f"test_{random_string(10)}@example.com"
        }

        cls.coffee_data_template = {
            "name": f"Test Coffee {random_string()}",
            "country": "Colombia", 
            "region": "Huila",    
            "farm": f"Test Farm {random_string()}",
            "variety": "Arabica",  
            "process": "Washed",   
            "roastProfile": "Medium", 
            "flavourNotes": ["citrus", "chocolate", "nuts"], 
            "description": f"Test coffee description {random_string(20)}",
            "imageUrl": f"https://example.com/{random_string()}.jpg",
        }
        cls.roastery_data_template_base = {
            "name": f"Test Roastery {random_string()}",
            "country": "Poland",         
            "city": "Katowice",            
            "address": "Korfantego 72",  
            "website": f"https://test{random_string()}.com/",
            "description": f"Test roastery description {random_string(20)}",
            "imageUrl": f"https://example.com/roastery_{random_string()}.jpg",
        }
        cls.shop_data_template_base = {
            "name": f"Test Shop {random_string()}",
            "country": "Poland",           
            "city": "Katowice",              
            "address": "Wawelska 1", 
            "website": f"https://testshop{random_string()}.com/",
            "description": f"Test shop description {random_string(20)}",
            "imageUrl": f"https://example.com/{random_string()}.jpg",
        }
        cls.review_payload_template_base = {
            "rating": random.randint(1, 5),
            "review": f"Test review {random_string(20)}",
        }

    def _get_roastery_data(self):
        data = self.roastery_data_template_base.copy()
        data["name"] = f"Test Roastery {random_string()}"
        data["website"] = f"https://test{random_string()}.com/"
        data["description"] = f"Test roastery description {random_string(20)}"
        data["imageUrl"] = f"https://example.com/roastery_{random_string()}.jpg"
        return data

    def _get_coffee_data(self):
        data = self.coffee_data_template.copy()
        data["name"] = f"Test Coffee {random_string()}"
        data["farm"] = f"Test Farm {random_string()}"
        data["description"] = f"Test coffee description {random_string(20)}"
        data["imageUrl"] = f"https://example.com/{random_string()}.jpg"
        return data

    def _get_shop_data(self):
        data = self.shop_data_template_base.copy()
        data["name"] = f"Test Shop {random_string()}"
        data["website"] = f"https://testshop{random_string()}.com/"
        data["description"] = f"Test shop description {random_string(20)}"
        data["imageUrl"] = f"https://example.com/shop_{random_string()}.jpg"
        return data

    def _get_review_payload(self):
        data = self.review_payload_template_base.copy()
        data["rating"] = random.randint(1, 5)
        data["review"] = f"Test review {random_string(20)}"
        return data

    def _create_roastery_for_test(self):
        payload = self._get_roastery_data()
        resp = requests.post(f"{BASE_URL}/roasteries", json=payload, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, f"Failed to create roastery for test: {resp.text} with payload {payload}")
        roastery = resp.json()
        self.assertIn("id", roastery)
        return roastery["id"]

    def _create_coffee_for_test(self, roastery_id):
        coffee_data = self._get_coffee_data()
        coffee_data["roasteryId"] = roastery_id
        resp = requests.post(f"{BASE_URL}/coffees", json=coffee_data, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, f"Failed to create coffee for test: {resp.text} with payload {coffee_data}")
        coffee = resp.json()
        self.assertIn("id", coffee)
        return coffee["id"]

    def _create_shop_for_test(self):
        payload = self._get_shop_data()
        resp = requests.post(f"{BASE_URL}/shops", json=payload, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, f"Failed to create shop for test: {resp.text} with payload {payload}")
        shop = resp.json()
        self.assertIn("id", shop)
        return shop["id"]

    def _create_review_for_test(self, review_target_payload):
        payload = self._get_review_payload()
        payload.update(review_target_payload)
        resp = requests.post(f"{BASE_URL}/reviews", json=payload, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, f"Failed to post review for test: {resp.text} with payload {payload}")
        review = resp.json()
        self.assertIn("id", review)
        return review["id"]

    def test_1_get_api_documentation_json(self):
        resp = requests.get(f"{BASE_URL}/")
        self.assertEqual(resp.status_code, 200)
        self.assertIn("application/json", resp.headers.get("Content-Type", ""))

    def test_2_get_api_documentation_html(self):
        resp = requests.get(f"{BASE_URL}/help")
        self.assertEqual(resp.status_code, 200)
        self.assertTrue(resp.text.strip().lower().startswith("<!doctype html>"))

    def test_3_get_user_by_id(self):
        resp = requests.get(f"{BASE_URL}/users/{self.user_id_for_tests}", headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, f"Failed to get user by ID: {resp.text}")
        body = resp.json()
        user = body.get("data") or body.get("user") or body 
        self.assertEqual(user.get("username"), self.admin_username)

    def test_4_create_roastery(self):
        roastery_id = self._create_roastery_for_test()
        self.assertIsNotNone(roastery_id)

    def test_5_get_all_roasteries(self):
        self._create_roastery_for_test()
        resp = requests.get(f"{BASE_URL}/roasteries")
        self.assertEqual(resp.status_code, 200, f"Failed to get roasteries: {resp.text}")
        data = resp.json()
        if data is None: 
            data = []
        self.assertIsInstance(data, list)
        self.assertGreater(len(data), 0)
        
    def test_6_get_roastery_by_id(self):
        roastery_id = self._create_roastery_for_test()
        resp = requests.get(f"{BASE_URL}/roasteries/{roastery_id}")
        self.assertEqual(resp.status_code, 200, f"Failed to get roastery by ID: {resp.text}")
        roastery = resp.json()
        self.assertIn("id", roastery)
        self.assertEqual(roastery["id"], roastery_id)
    
    def test_7_update_roastery(self):
        roastery_id = self._create_roastery_for_test()
        
        update_payload = self._get_roastery_data()
        update_payload["name"] = f"Updated Roastery {random_string()}"
        update_payload["description"] = f"Updated description {random_string(20)}"
        if "id" in update_payload: 
            del update_payload["id"]

        resp = requests.put(f"{BASE_URL}/roasteries/{roastery_id}", json=update_payload, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, f"Failed to update roastery: {resp.text} with payload {update_payload}")
        updated_roastery = resp.json()
        self.assertEqual(updated_roastery["name"], update_payload["name"])
        self.assertEqual(updated_roastery["description"], update_payload["description"])


    def test_8_create_coffee(self):
        roastery_id = self._create_roastery_for_test()
        coffee_id = self._create_coffee_for_test(roastery_id)
        self.assertIsNotNone(coffee_id)

    def test_9_get_all_coffees(self):
        roastery_id = self._create_roastery_for_test()
        self._create_coffee_for_test(roastery_id)
        resp = requests.get(f"{BASE_URL}/coffees")
        self.assertEqual(resp.status_code, 200, f"Failed to get coffees: {resp.text}")
        data = resp.json()
        if data is None: 
            data = []
        self.assertIsInstance(data, list)
        self.assertGreater(len(data), 0)
        
    def test_10_get_coffee_by_id(self):
        roastery_id = self._create_roastery_for_test()
        coffee_id = self._create_coffee_for_test(roastery_id)
        resp = requests.get(f"{BASE_URL}/coffees/{coffee_id}")
        self.assertEqual(resp.status_code, 200, f"Failed to get coffee by ID: {resp.text}")
        coffee = resp.json()
        self.assertIn("id", coffee)
        self.assertEqual(coffee["id"], coffee_id)
        
    def test_11_update_coffee(self):
        roastery_id = self._create_roastery_for_test()
        coffee_id = self._create_coffee_for_test(roastery_id)
        
        get_coffee_resp = requests.get(f"{BASE_URL}/coffees/{coffee_id}")
        self.assertEqual(get_coffee_resp.status_code, 200)
        current_coffee_data = get_coffee_resp.json()

        update_payload = current_coffee_data
        update_payload["name"] = f"Updated Coffee {random_string()}"
        update_payload["description"] = f"Updated description {random_string(20)}"
        if "id" in update_payload:
             del update_payload["id"]
        if "roastery" in update_payload: # If API returns nested roastery, remove for PUT
            del update_payload["roastery"]


        resp = requests.put(f"{BASE_URL}/coffees/{coffee_id}", json=update_payload, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, f"Failed to update coffee: {resp.text} with payload {update_payload}")
        updated_coffee = resp.json()
        self.assertEqual(updated_coffee["name"], update_payload["name"])
        self.assertEqual(updated_coffee["description"], update_payload["description"])


    def test_12_create_coffee_shop(self):
        shop_id = self._create_shop_for_test()
        self.assertIsNotNone(shop_id)
        
    def test_13_get_all_coffee_shops(self):
        self._create_shop_for_test()
        resp = requests.get(f"{BASE_URL}/shops")
        self.assertEqual(resp.status_code, 200, f"Failed to get coffee shops: {resp.text}")
        data = resp.json()
        if data is None: 
            data = []
        self.assertIsInstance(data, list)
        self.assertGreater(len(data), 0)

    def test_14_get_coffee_shop_by_id(self):
        shop_id = self._create_shop_for_test()
        resp = requests.get(f"{BASE_URL}/shops/{shop_id}")
        self.assertEqual(resp.status_code, 200, f"Failed to get coffee shop by ID: {resp.text}")
        shop = resp.json()
        self.assertIn("id", shop)
        self.assertEqual(shop["id"], shop_id)

    def test_15_update_coffee_shop(self):
        shop_id = self._create_shop_for_test()

        update_payload = self._get_shop_data()
        update_payload["name"] = f"Updated Coffee Shop {random_string()}"
        update_payload["description"] = f"Updated description {random_string(20)}"
        if "id" in update_payload: 
            del update_payload["id"]
        
        resp = requests.put(f"{BASE_URL}/shops/{shop_id}", json=update_payload, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, f"Failed to update coffee shop: {resp.text} with payload {update_payload}")
        updated_shop = resp.json()
        self.assertEqual(updated_shop["name"], update_payload["name"])
        self.assertEqual(updated_shop["description"], update_payload["description"])


    def test_16_create_coffee_review(self):
        roastery_id = self._create_roastery_for_test()
        coffee_id = self._create_coffee_for_test(roastery_id)
        review_id = self._create_review_for_test({"coffeeId": coffee_id})
        self.assertIsNotNone(review_id)

    def test_17_create_roastery_review(self):
        roastery_id = self._create_roastery_for_test()
        review_id = self._create_review_for_test({"roasteryId": roastery_id})
        self.assertIsNotNone(review_id)

    def test_18_create_coffee_shop_review(self):
        shop_id = self._create_shop_for_test()
        review_id = self._create_review_for_test({"coffeeShopId": shop_id})
        self.assertIsNotNone(review_id)

    def test_19_get_all_reviews(self):
        roastery_id = self._create_roastery_for_test()
        coffee_id = self._create_coffee_for_test(roastery_id)
        self._create_review_for_test({"coffeeId": coffee_id})
        resp = requests.get(f"{BASE_URL}/reviews")
        self.assertEqual(resp.status_code, 200, f"Failed to get all reviews: {resp.text}")
        reviews = resp.json()
        if reviews is None: 
            reviews = []
        self.assertIsInstance(reviews, list)
        self.assertGreater(len(reviews),0)
        
    def test_20_get_reviews_for_coffee(self):
        roastery_id = self._create_roastery_for_test()
        coffee_id = self._create_coffee_for_test(roastery_id)
        self._create_review_for_test({"coffeeId": coffee_id}) 
        resp = requests.get(f"{BASE_URL}/reviews?coffeeId={coffee_id}")
        self.assertEqual(resp.status_code, 200, f"Failed to get coffee reviews: {resp.text}")
        reviews = resp.json()
        if reviews is None: 
            reviews = []
        self.assertIsInstance(reviews, list)
        self.assertGreater(len(reviews),0)
        for review in reviews:
            self.assertEqual(str(review.get("coffeeId") or review.get("coffee", {}).get("id")), str(coffee_id))

    def test_21_get_reviews_for_roastery(self):
        roastery_id = self._create_roastery_for_test()
        self._create_review_for_test({"roasteryId": roastery_id}) 
        resp = requests.get(f"{BASE_URL}/reviews?roasteryId={roastery_id}")
        self.assertEqual(resp.status_code, 200, f"Failed to get roastery reviews: {resp.text}")
        reviews = resp.json()
        if reviews is None: 
            reviews = []
        self.assertIsInstance(reviews, list)
        self.assertGreater(len(reviews),0)
        for review in reviews:
            self.assertEqual(str(review.get("roasteryId") or review.get("roastery", {}).get("id")), str(roastery_id))

    def test_22_get_reviews_for_shop(self):
        shop_id = self._create_shop_for_test()
        self._create_review_for_test({"coffeeShopId": shop_id}) 
        resp = requests.get(f"{BASE_URL}/reviews?coffeeShopId={shop_id}")
        self.assertEqual(resp.status_code, 200, f"Failed to get shop reviews: {resp.text}")
        reviews = resp.json()
        if reviews is None: 
            reviews = []
        self.assertIsInstance(reviews, list)
        self.assertGreater(len(reviews),0)
        for review in reviews:
            self.assertEqual(str(review.get("coffeeShopId") or review.get("coffeeShop", {}).get("id")), str(shop_id))

    def test_23_get_reviews_for_user(self):
        roastery_id = self._create_roastery_for_test()
        self._create_review_for_test({"roasteryId": roastery_id}) 
        resp = requests.get(f"{BASE_URL}/reviews?userId={self.user_id_for_tests}")
        self.assertEqual(resp.status_code, 200, f"Failed to get user reviews: {resp.text}")
        reviews = resp.json()
        if reviews is None: 
            reviews = []
        self.assertIsInstance(reviews, list)
        self.assertGreater(len(reviews),0)
        for review in reviews:
            self.assertEqual(str(review.get("userId") or review.get("user", {}).get("id")), str(self.user_id_for_tests))

    def test_24_update_review(self):
        roastery_id = self._create_roastery_for_test()
        coffee_id = self._create_coffee_for_test(roastery_id)
        review_id = self._create_review_for_test({"coffeeId": coffee_id})
        
        update_payload = {
            "rating": random.randint(1, 5),
            "review": f"Updated review text {random_string()}"
        }
        resp = requests.put(f"{BASE_URL}/reviews/{review_id}", json=update_payload, headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, f"Failed to update review: {resp.text}")
        updated_review = resp.json()
        self.assertEqual(updated_review["rating"], update_payload["rating"])
        self.assertEqual(updated_review["review"], update_payload["review"])
        self.assertEqual(updated_review["id"], review_id)

    def test_25_get_stats(self):
        resp = requests.get(f"{BASE_URL}/stats", headers=self.auth_headers)
        self.assertEqual(resp.status_code, 200, f"Failed to get stats: {resp.text}")
        stats = resp.json()
        self.assertIsInstance(stats, dict)

    def test_26_delete_coffee_review(self):
        roastery_id = self._create_roastery_for_test()
        coffee_id = self._create_coffee_for_test(roastery_id)
        review_id = self._create_review_for_test({"coffeeId": coffee_id})
        
        resp_delete = requests.delete(f"{BASE_URL}/reviews/{review_id}", headers=self.auth_headers)
        self.assertEqual(resp_delete.status_code, 204, f"Failed to delete coffee review: {resp_delete.text}")

    def test_27_delete_roastery_review(self):
        roastery_id = self._create_roastery_for_test()
        review_id = self._create_review_for_test({"roasteryId": roastery_id})
        
        resp_delete = requests.delete(f"{BASE_URL}/reviews/{review_id}", headers=self.auth_headers)
        self.assertEqual(resp_delete.status_code, 204, f"Failed to delete roastery review: {resp_delete.text}")
        
    def test_28_delete_coffee_shop_review(self):
        shop_id = self._create_shop_for_test()
        review_id = self._create_review_for_test({"coffeeShopId": shop_id})
        
        resp_delete = requests.delete(f"{BASE_URL}/reviews/{review_id}", headers=self.auth_headers)
        self.assertEqual(resp_delete.status_code, 204, f"Failed to delete coffee shop review: {resp_delete.text}")
        
    def test_29_delete_coffee(self):
        roastery_id = self._create_roastery_for_test()
        coffee_id = self._create_coffee_for_test(roastery_id)
        
        resp_delete = requests.delete(f"{BASE_URL}/coffees/{coffee_id}", headers=self.auth_headers)
        self.assertEqual(resp_delete.status_code, 204, f"Failed to delete coffee: {resp_delete.text}")

        resp_get = requests.get(f"{BASE_URL}/coffees/{coffee_id}")
        self.assertEqual(resp_get.status_code, 404)


    def test_30_delete_roastery(self):
        roastery_id = self._create_roastery_for_test()
        
        resp_delete = requests.delete(f"{BASE_URL}/roasteries/{roastery_id}", headers=self.auth_headers)
        self.assertEqual(resp_delete.status_code, 204, f"Failed to delete roastery: {resp_delete.text}")
        
        resp_get = requests.get(f"{BASE_URL}/roasteries/{roastery_id}")
        self.assertEqual(resp_get.status_code, 404)
        
    def test_31_delete_coffee_shop(self): 
        shop_id = self._create_shop_for_test()
        
        resp_delete = requests.delete(f"{BASE_URL}/shops/{shop_id}", headers=self.auth_headers)
        self.assertEqual(resp_delete.status_code, 204, f"Failed to delete coffee shop: {resp_delete.text}")

        resp_get = requests.get(f"{BASE_URL}/shops/{shop_id}")
        self.assertEqual(resp_get.status_code, 404)


    def test_32_register_and_login_new_user(self):
        user_data = self.new_user_register_data_template.copy()
        user_data["username"] = f"testuser_{random_string(10)}"
        user_data["email"] = f"test_{random_string(10)}@example.com"

        reg_resp = requests.post(f"{BASE_URL}/register", json=user_data)
        self.assertEqual(reg_resp.status_code, 200, f"New user registration failed: {reg_resp.text} with payload {user_data}")
        
        user_json = reg_resp.json()
        new_user_id_val = user_json.get("id") or user_json.get("user_id") or user_json.get("userId")
        if not new_user_id_val and "user" in user_json and isinstance(user_json.get("user"), dict):
            new_user_id_val = user_json["user"].get("id")
        self.assertIsNotNone(new_user_id_val, f"New user ID not found in registration response: {user_json}")

        time.sleep(1) 

        new_user_login_payload = {
            "username": user_data["username"],
            "passwords": user_data["password"] 
        }
        
        login_resp = requests.post(f"{BASE_URL}/login", json=new_user_login_payload)
        self.assertEqual(login_resp.status_code, 200, f"New user login failed: {login_resp.text} with payload {new_user_login_payload}")
        
        token_data = login_resp.json()
        new_user_token = token_data.get("token") or token_data.get("access_token")
        self.assertIsNotNone(new_user_token, f"No token returned from new user login: {token_data}")
        logging.info(f"Successfully registered and logged in new user: {user_data['username']}")

if __name__ == "__main__":
    unittest.main()