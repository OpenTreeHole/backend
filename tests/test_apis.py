from rest_framework.test import APITestCase


class IndexTests(APITestCase):
    """hi"""

    def test_get(self):
        r = self.client.get("/")
        self.assertEqual(r.data, {"message": "Hello world!"})
