package com.seanglay.sombot_easy.model;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;

import java.time.LocalDateTime;

@Entity
@Table(name = "promotions")
@Data
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode
public class Promotion {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "promotion_id")
    private Long promotionId;

    @NotNull(message = "Promotion name cannot be null")
    @Size(min = 3, max = 255, message = "Promotion name must be between 3 and 255 characters")
    @Column(name = "promotion_name", nullable = false)
    private String promotionName;

    @NotNull(message = "Image URL cannot be null")
    @Size(min = 5, max = 500, message = "Image URL must be between 5 and 500 characters")
    @Column(name = "image_url", nullable = false)
    private String imageUrl;

    @Column(name = "start_date", nullable = false)
    private LocalDateTime startDate;

    @Column(name = "end_date", nullable = false)
    private LocalDateTime endDate;

    @Enumerated(EnumType.STRING)
    @Column(name = "status", columnDefinition = "varchar default 'active'")
    private PromotionStatus status = PromotionStatus.ACTIVE;

    @Column(name = "created_at", updatable = false, nullable = false)
    private LocalDateTime createdAt;

    @Column(name = "updated_at")
    private LocalDateTime updatedAt;

    @PrePersist
    public void prePersist() {
        LocalDateTime now = LocalDateTime.now();
        createdAt = now;
        updatedAt = now;
    }

    @PreUpdate
    public void preUpdate() {
        updatedAt = LocalDateTime.now();
    }

    public enum PromotionStatus {
        ACTIVE, INACTIVE, EXPIRED
    }
}
