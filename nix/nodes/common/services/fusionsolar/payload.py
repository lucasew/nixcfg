#!/usr/bin/env -S sd nix shell --really
#!nix-shell -i python3 -p python3Packages.selenium chromedriver chromium

from urllib.request import urlopen, Request
import json
from argparse import ArgumentParser
from pathlib import Path
import time
import base64
import os

from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By

parser = ArgumentParser()
parser.add_argument("--user", os.getenv("FUSIONSOLAR_USER"))
parser.add_argument("--password", os.getenv("FUSIONSOLAR_PASSWORD"))
parser.add_argument("--output", type=Path, required=True)
args = parser.parse_args()

driver = webdriver.Chrome()
print('[*] Login')
driver.get("https://intl.fusionsolar.huawei.com/pvmswebsite/login/build/index.html#/LOGIN")
driver.find_element(By.CSS_SELECTOR, "div#username > input").send_keys(args.user)
password_input = driver.find_element(By.CSS_SELECTOR, "div#password > input")
password_input.send_keys(args.password)
password_input.send_keys(Keys.ENTER)
time.sleep(10)
print('[*] Homepage')
driver.get("https://intl.fusionsolar.huawei.com")
time.sleep(10)
print('[*] Tentando listar estações')
stations = driver.find_elements(By.CSS_SELECTOR, "tbody.ant-table-tbody a.nco-home-list-text-ellipsis")
print('stations', stations)

stations_data = []
for station in stations:
    station_data = dict(
        url=station.get_attribute('href'),
        name=station.text 
    )
    print(station_data)
    stations_data.append(station_data)

for station in stations_data:
    station_url = station['url']
    station_name = station['name']
    print(f'[*] Chupinhando dados da estação "{station_name}"')
    driver.get(station_url)
    time.sleep(10)
    the_canvas = driver.find_element(By.CSS_SELECTOR, ".nco-single-energy-body canvas")
    canvas_b64 = driver.execute_script("return arguments[0].toDataURL('image/png').substring(21);", the_canvas)
    amount_produced = float(driver.find_element(By.CSS_SELECTOR, "span.value").text)
    print(f'[*] Produzido hoje: {amount_produced}kWh')
    print(f'[*] Salvando dados da estação "{station_name}"')
    (args.output / f"{station_name}.png").write_bytes(base64.b64decode(canvas_b64))
    
# driver.sleep(600)
time.sleep(60)

print(driver.title)
# token_file = Path.home() / ".fusion-solar"

# token = None
# if token_file.exists():
#     token = token_file.read_text()
# else:
#     request = Request(
#         "https://intl.fusionsolar.huawei.com/thirdData/login",
#         data=json.dumps(dict(
#             userName = args.user,
#             systemCode = args.password
#         )).encode('utf-8'),
#         headers={
#             "accept": "application/json",
#             "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
#         },
#         method='POST'
#     )
#     data = urlopen(request)
#     token = data.getheader("xsrf-token")
#     if token is not None:
#         token_file.write_text(token)
#     print(data)
#     print(data.read())
