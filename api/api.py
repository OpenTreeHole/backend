from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.views import APIView


# class Index(APIView):
@api_view()
def index(request):
    return Response({"message": "Hello world!"})
