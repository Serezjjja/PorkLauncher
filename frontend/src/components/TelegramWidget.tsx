import React from "react";
import { useTranslation } from "../i18n";

const TELEGRAM_LINK = "https://t.me/porklandmc";
const TG_LOGO_URL = "https://porkland.net/storage/img/tg.svg";

export const TelegramWidget: React.FC = () => {
  const { language } = useTranslation();

  const getTitleText = () => {
    return language === "ru" ? "Наш Telegram" : "Our Telegram";
  };

  const getDescriptionText = () => {
    return language === "ru" ? "Новости и обновления" : "News and updates";
  };

  const getButtonText = () => {
    return language === "ru" ? "ПРИСОЕДИНИТЬСЯ" : "JOIN";
  };

  return (
    <div className="telegram-widget">
      {/* Background Telegram Icon */}
      <img src={TG_LOGO_URL} className="telegram-bg-icon" alt="" />

      {/* Header */}
      <div className="tg-header">
        <div className="tg-logo-wrap">
          <img src={TG_LOGO_URL} alt="TG Logo" />
        </div>
        <div className="tg-info">
          <h3>{getTitleText()}</h3>
          <p>{getDescriptionText()}</p>
        </div>
      </div>

      {/* Join Button */}
      <a
        href={TELEGRAM_LINK}
        target="_blank"
        rel="noopener noreferrer"
        className="tg-btn"
      >
        {getButtonText()}
      </a>
    </div>
  );
};
