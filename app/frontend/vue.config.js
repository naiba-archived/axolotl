const proxyConfig = {
  target: "http://localhost",
  ws: true,
  changeOrigin: true,
  // secure: true,
  onProxyRes: function(proxyRes, req) {
    const cookies = proxyRes.headers["set-cookie"];
    const cookieRegex = /Domain=localhost/i;
    if (cookies) {
      const newCookie = cookies.map(function(cookie) {
        if (cookieRegex.test(cookie)) {
          return cookie.replace(
            cookieRegex,
            "Domain=" + req.headers.host.split(":")[0]
          );
        }
        return cookie;
      });
      delete proxyRes.headers["set-cookie"];
      newCookie[0] = newCookie[0].replace("; Secure; SameSite=None", "");
      proxyRes.headers["set-cookie"] = newCookie;
    }
  }
};

module.exports = {
  pluginOptions: {
    vconsole: { enable: false }
  },
  lintOnSave: false,
  devServer: {
    disableHostCheck: true,
    proxy: {
      "^/api": proxyConfig
    }
  }
};
