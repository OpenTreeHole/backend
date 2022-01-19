from api.urls import urlpatterns as http_urlpatterns
from ws.urls import urlpatterns as ws_urlpatterns

urlpatterns = http_urlpatterns + ws_urlpatterns
