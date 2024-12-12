package com.seanglay.sombot_easy.repository;

import com.seanglay.sombot_easy.model.Promotion;
import org.springframework.data.jpa.repository.JpaRepository;

public interface PromotionRepository extends JpaRepository<Promotion, Long> {
}
