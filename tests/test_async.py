from rest_framework.test import APITestCase


class IndexTests(APITestCase):
    """hi"""

    def test_get(self):
        r = self.client.get("/")
        print(r.data)
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.data, {"message": "Hello world!"})
