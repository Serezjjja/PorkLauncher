import React, { useEffect, useState } from "react";
import { useTranslation } from "../i18n";

interface DiscordMember {
  id: string;
  username: string;
  avatar_url: string;
}

interface DiscordData {
  presence_count: number;
  instant_invite: string;
  members: DiscordMember[];
}

const SERVER_ID = "1227227474551767101";
const DEFAULT_INVITE = "https://discord.gg/RbreKRwsH7";

export const DiscordWidget: React.FC = () => {
  const { t, language } = useTranslation();
  const [discordData, setDiscordData] = useState<DiscordData | null>(null);
  const [error, setError] = useState(false);

  useEffect(() => {
    const fetchDiscordData = async () => {
      try {
        const response = await fetch(
          `https://discord.com/api/guilds/${SERVER_ID}/widget.json`
        );
        if (!response.ok) throw new Error("Failed to fetch");
        const data = await response.json();
        setDiscordData(data);
        setError(false);
      } catch (err) {
        console.log("Discord widget error:", err);
        setError(true);
      }
    };

    fetchDiscordData();
    // Refresh every 60 seconds
    const interval = setInterval(fetchDiscordData, 60000);
    return () => clearInterval(interval);
  }, []);

  const onlineCount = discordData?.presence_count || 0;
  const inviteUrl = discordData?.instant_invite || DEFAULT_INVITE;
  const members = discordData?.members?.slice(0, 5) || [];
  const moreCount = Math.max(0, onlineCount - 5);

  const getOnlineText = () => {
    if (error) return language === "ru" ? "Discord Сервер" : "Discord Server";
    if (!discordData) return language === "ru" ? "Загрузка..." : "Loading...";
    return language === "ru" ? `${onlineCount} в сети` : `${onlineCount} online`;
  };

  const getButtonText = () => {
    return language === "ru" ? "ПРИСОЕДИНИТЬСЯ" : "JOIN";
  };

  const getTitleText = () => {
    return language === "ru" ? "Наш Discord" : "Our Discord";
  };

  return (
    <div className="discord-widget">
      {/* Background Discord Icon */}
      <svg className="discord-bg-icon" viewBox="0 0 24 24">
        <path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a16.09 16.09 0 0 0-4.8 0c-.14-.34-.35-.76-.54-1.09-.01-.02-.04-.03-.07-.03c-1.5.26-2.93.71-4.27 1.33c-.01 0-.02.01-.03.02c-2.72 4.07-3.47 8.03-3.1 11.95c0 .02.01.04.03.05c1.8 1.32 3.53 2.12 5.2 2.65c.03.01.06 0 .07-.02c.4-.55.76-1.13 1.07-1.74c.02-.04 0-.08-.04-.09c-.57-.22-1.11-.48-1.64-.78c-.04-.02-.04-.08.01-.11c.11-.08.22-.17.33-.25c.02-.02.05-.02.07-.01c3.44 1.57 7.15 1.57 10.55 0c.02-.01.05-.01.07.01c.11.09.22.17.33.26c.04.03.04.09 0 .11a10.9 10.9 0 0 1-1.64.78c-.04.01-.05.06-.04.09c.31.61.66 1.19 1.07 1.74c.03.01.06.02.09.01c1.72-.53 3.45-1.33 5.25-2.65c.02-.01.03-.03.03-.05c.44-4.53-.73-8.46-3.1-11.95c-.01-.01-.02-.02-.04-.02zM8.52 14.91c-1.03 0-1.89-.95-1.89-2.12s.84-2.12 1.89-2.12c1.06 0 1.9.96 1.89 2.12c0 1.17-.84 2.12-1.89 2.12zm6.97 0c-1.03 0-1.89-.95-1.89-2.12s.84-2.12 1.89-2.12c1.06 0 1.9.96 1.89 2.12c0 1.17-.85 2.12-1.89 2.12z" />
      </svg>

      {/* Header */}
      <div className="discord-header">
        <div className="discord-logo-wrap">
          <svg viewBox="0 0 24 24">
            <path d="M19.27 5.33C17.94 4.71 16.5 4.26 15 4a.09.09 0 0 0-.07.03c-.18.33-.39.76-.53 1.09a16.09 16.09 0 0 0-4.8 0c-.14-.34-.35-.76-.54-1.09c-.01-.02-.04-.03-.07-.03c-1.5.26-2.93.71-4.27 1.33c-.01 0-.02.01-.03.02c-2.72 4.07-3.47 8.03-3.1 11.95c0 .02.01.04.03.05c1.8 1.32 3.53 2.12 5.2 2.65c.03.01.06 0 .07-.02c.4-.55.76-1.13 1.07-1.74c.02-.04 0-.08-.04-.09c-.57-.22-1.11-.48-1.64-.78c-.04-.02-.04-.08.01-.11c.11-.08.22-.17.33-.25c.02-.02.05-.02.07-.01c3.44 1.57 7.15 1.57 10.55 0c.02-.01.05-.01.07.01c.11.09.22.17.33.26c.04.03.04.09 0 .11a10.9 10.9 0 0 1-1.64.78c-.04.01-.05.06-.04.09c.31.61.66 1.19 1.07 1.74c.03.01.06.02.09.01c1.72-.53 3.45-1.33 5.25-2.65c.02-.01.03-.03.03-.05c.44-4.53-.73-8.46-3.1-11.95c-.01-.01-.02-.02-.04-.02zM8.52 14.91c-1.03 0-1.89-.95-1.89-2.12s.84-2.12 1.89-2.12c1.06 0 1.9.96 1.89 2.12c0 1.17-.84 2.12-1.89 2.12zm6.97 0c-1.03 0-1.89-.95-1.89-2.12s.84-2.12 1.89-2.12c1.06 0 1.9.96 1.89 2.12c0 1.17-.85 2.12-1.89 2.12z" />
          </svg>
        </div>
        <div className="discord-info">
          <h3>{getTitleText()}</h3>
          <p>
            <span className="online-dot"></span>
            <span>{getOnlineText()}</span>
          </p>
        </div>
      </div>

      {/* Members List */}
      <div className="members-list">
        {members.map((member, index) => (
          <div
            key={member.id}
            className="member-avatar"
            style={{
              zIndex: 5 - index,
              backgroundImage: `url('${member.avatar_url}')`,
            }}
            title={member.username}
          />
        ))}
        {moreCount > 0 && <div className="more-members">+{moreCount}</div>}
      </div>

      {/* Join Button */}
      <a
        href={inviteUrl}
        target="_blank"
        rel="noopener noreferrer"
        className="discord-btn"
      >
        {getButtonText()}
      </a>
    </div>
  );
};
