package com.seanglay.sombot_easy.util;

import com.seanglay.sombot_easy.model.User;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.io.Decoders;
import io.jsonwebtoken.security.Keys;
import org.springframework.stereotype.Component;

import java.security.Key;
import java.sql.Date;
import java.util.HashMap;
import java.util.Map;

@Component
public class JWTToken {

    public String generateToken(User user) {

        Map<String, Object> claims = new HashMap<>();
        claims.put("id", user.getUserId());
        claims.put("username", user.getUsername());

        return Jwts.builder().claims().add(claims).subject(user.getEmail()).issuedAt(new Date(System.currentTimeMillis())).expiration(new Date(System.currentTimeMillis() + 60 * 60 * 60 * 10)).and().signWith(getKey()).compact();
    }

    private Key getKey() {
        String secretKey = "xsEm_mpWkqlXpv7O6TUVLZnHOKyvZ6v_Req7w9Day3U";
        byte[] keyBytes = Decoders.BASE64URL.decode(secretKey);
        return Keys.hmacShaKeyFor(keyBytes);
    }
}
