import React from "react";

export type BannerVariant = "large" | "small" | "compact";

export interface BannerProps {
  variant: BannerVariant;
  backgroundImage?: string;
  iconImage?: string;
  title?: string;
  description?: string;
  text?: string;
  position?: {
    left?: string;
    top?: string;
  };
  onClick?: () => void;
  className?: string;
}

function Banner({
  variant,
  backgroundImage,
  iconImage,
  title,
  description,
  text,
  position,
  onClick,
  className = "",
}: BannerProps) {
  // Base classes common to all variants
  const baseClasses = "rounded-[20px] border border-[#5D8B57]/[0.25] overflow-hidden group cursor-pointer transform-gpu shadow-[0_12px_48px_rgba(0,0,0,0.4),0_0_30px_rgba(93,139,87,0.1),inset_0_1px_0_rgba(255,255,255,0.1)] hover:shadow-[0_12px_48px_rgba(0,0,0,0.4),0_0_40px_rgba(93,139,87,0.2),inset_0_1px_0_rgba(255,255,255,0.15)] transition-all duration-300";

  // Variant-specific styling
  const variantClasses = {
    large: "w-[448px] h-[200px]",
    small: "w-[213px] h-[200px] bg-[#364E3A]/[0.65] backdrop-blur-[24px]",
    compact: "w-[400px] h-[80px] border-[#5D8B57]/[0.3] bg-[#364E3A]/65 backdrop-blur-[24px] px-[10px]",
  };

  // Position classes - use absolute if position is provided
  const positionClasses = position ? "absolute" : "";

  const containerClasses = `${baseClasses} ${variantClasses[variant]} ${positionClasses} ${className}`;

  // Inline styles for dynamic positioning
  const positionStyle: React.CSSProperties = position
    ? {
        left: position.left,
        top: position.top,
      }
    : {};

  // Large variant: Full featured banner with background, overlay, icon, title, and description
  if (variant === "large") {
    return (
      <button className={containerClasses} style={positionStyle} onClick={onClick}>
        {/* Background image */}
        {backgroundImage && (
          <img
            src={backgroundImage}
            alt=""
            className="w-full h-full object-cover opacity-90 transition-all duration-300 filter saturate-[0.6] contrast-[0.85] brightness-[0.93] group-hover:saturate-100 group-hover:contrast-100 group-hover:brightness-100 will-change-[filter]"
          />
        )}

        {/* Dark overlay */}
        <div className="absolute inset-0 bg-[#1C2A1F]/30 pointer-events-none" />

        {/* Small icon */}
        {iconImage && (
          <img
            src={iconImage}
            alt="Banner icon"
            className="absolute bottom-[10px] left-[10px] w-[60px] h-[60px] rounded-[10px] pointer-events-none transform-gpu"
          />
        )}

        {/* Text block */}
        {(title || description) && (
          <div
            className="
              absolute bottom-[14px]
              left-[80px]
              right-[14px]
              w-[310px]
              flex flex-col
              pointer-events-none
            "
          >
            {/* Title */}
            {title && (
              <div className="font-[Unbounded] text-[14px] text-white/90 text-left drop-shadow-[0_0_8px_rgba(93,139,87,0.5)]">
                {title}
              </div>
            )}

            {/* Description */}
            {description && (
              <p className="mt-[2px] text-[14px] leading-[16px] font-[Mazzard] text-white/85 text-justify drop-shadow-[0_0_4px_rgba(0,0,0,0.5)]">
                {description}
              </p>
            )}
          </div>
        )}
      </button>
    );
  }

  // Small variant: Simple text-only banner
  if (variant === "small") {
    return (
      <button className={containerClasses} style={positionStyle} onClick={onClick}>
        {text && (
          <p className="text-[14px] font-[Unbounded] text-white/60 drop-shadow-[0_0_8px_rgba(93,139,87,0.3)]">
            {text}
          </p>
        )}
      </button>
    );
  }

  // Compact variant: Horizontal layout with image and text side by side
  if (variant === "compact") {
    return (
      <div className={`${containerClasses} flex items-center gap-[12px]`} style={positionStyle} onClick={onClick}>
        {/* Image */}
        {iconImage && (
          <img
            src={iconImage}
            alt="Banner icon"
            className="w-[60px] h-[60px] rounded-[10px]"
          />
        )}

        {/* Text with title and description */}
        {(title || description) ? (
          <div className="flex flex-col justify-center overflow-hidden">
            {title && (
              <span className="text-[14px] font-[Mazzard] text-white/90 tracking-[-3%] drop-shadow-[0_0_6px_rgba(93,139,87,0.4)]">
                {title}
              </span>
            )}
            {description && (
              <span className="text-[12px] font-[Mazzard] text-white/50 tracking-[-3%] truncate max-w-[300px]">
                {description}
              </span>
            )}
          </div>
        ) : text && (
          <div className="flex flex-col justify-center">
            <span className="text-[14px] text-center text-white/90 font-[Mazzard] tracking-[-3%] drop-shadow-[0_0_6px_rgba(93,139,87,0.4)]">
              {text}
            </span>
          </div>
        )}
      </div>
    );
  }

  return null;
}

export default Banner;

