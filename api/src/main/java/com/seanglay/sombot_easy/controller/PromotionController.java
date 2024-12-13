package com.seanglay.sombot_easy.controller;

import com.seanglay.sombot_easy.dto.APIResponse;
import com.seanglay.sombot_easy.model.Promotion;
import com.seanglay.sombot_easy.service.PromotionService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;

@RestController
@RequestMapping("/api/promotions")
public class PromotionController extends BaseController {

    private final PromotionService promotionService;

    @Autowired
    public PromotionController(PromotionService promotionService) {
        this.promotionService = promotionService;
    }

    @PostMapping
    public ResponseEntity<APIResponse<Promotion>> createPromotion(@RequestBody Promotion promotion) {
        Promotion createdPromotion = promotionService.createPromotion(promotion);
        return ok(createdPromotion);
    }

    @GetMapping
    public ResponseEntity<APIResponse<List<Promotion>>> getAllPromotions() {
        List<Promotion> promotions = promotionService.getAllPromotions();
        return ok(promotions);
    }

    @GetMapping("/{id}")
    public ResponseEntity<APIResponse<Promotion>> getPromotionById(@PathVariable Long id) {
        Optional<Promotion> promotion = promotionService.getPromotionById(id);
        return promotion.map(this::ok).orElseGet(() -> ResponseEntity.notFound().build());
    }

    @PutMapping("/{id}")
    public ResponseEntity<APIResponse<Promotion>> updatePromotion(@PathVariable Long id, @RequestBody Promotion promotion) {
        Promotion updatedPromotion = promotionService.updatePromotion(id, promotion);
        return ok(updatedPromotion);
    }

    @DeleteMapping("/{id}")
    public boolean deletePromotion(@PathVariable Long id) {
        return promotionService.deletePromotion(id);
    }
}
