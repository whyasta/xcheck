wrk.headers["Cookie"] = "_gorilla_csrf=MTcyODI4MDYyM3xJaTlqUldsaEszcEtjVTFwU25KQ1NWRnJSa3hrZVZGbGRqWnhTazVpWVcxR2JURlpSa0Y1WWxwNmMxazlJZ289fGRLL3W45GY4fNXgt_lEfekH8CyWJE3FiJiFRoe6BaVJ; accessToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiI3Y2JkZjZhMy0yY2I2LTQ0ZGItODE2Yy1hNWQxODJhZDVjODUiLCJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3Mjg0MDU3MDAsInN1YiI6MSwidXNlcl9pZCI6MX0.JxaT0XwY4RVfmeLuPng5HattBlsFGSvvU0QGav5yME8; refreshToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg0MDU3MDAsInN1YiI6MX0.kd6rG7FPawo8ic4SNWWtyCozEkMy50wWfV1TIifceQA"
response = function(status, headers, body)
    io.write("status:" .. status .. "\n")
    -- io.write("status:" .. status .. "\n" .. body .. "\n-------------------------------------------------\n")
    -- for key, value in pairs(headers) do
    --     if key == "Location" then
    --         io.write("Location header found!\n")
    --         io.write(key)
    --         io.write(":")
    --         io.write(value)
    --         io.write("\n")
    --         io.write("---\n")
    --         break
    --     end
    -- end
end
