local wrk = require "wrk"

function randomInt(min, max)
    return math.random(min, max)
end

function randomString(length)
    local result = ""
    local caracteres = "            abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
      for i = 1, length do
         local index = math.random(1, string.len(caracteres))
         result = result .. string.sub(caracteres, index, index)
      end
      return result
end

local contentID = randomInt(1, 10)
local rating = randomInt(1, 10)
local text = randomString(100)
local title = randomString(20)

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["x-csrf"] = "vUMyVabJPtByzqqXteXivqHmHuZGcvfs"
wrk.headers["Cookie"] = "_csrf=vUMyVabJPtByzqqXteXivqHmHuZGcvfs; session=392fbbef-aa18-4149-9724-d3d6fa056d50"
wrk.body = '{"contentID": ' .. tostring(contentID) .. ',"rating": ' .. tostring(rating) .. ',"text": "' .. tostring(text) .. '","title": "' .. tostring(title) .. '"}'

--- при желании можно вывести в файл или консоль результат запроса

response = function(status, headers, body)
  print("Status: " .. status)
  print("Body: " .. body)
end

