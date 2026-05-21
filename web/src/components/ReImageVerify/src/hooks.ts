import { ref, onMounted } from "vue";
import { getCaptcha } from "@/api/user";
import { message } from "@/utils/message";
import { debounce } from "@pureadmin/utils";

function escapeSvg(text: string) {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/\"/g, "&quot;")
    .replace(/'/g, "&apos;");
}

function buildCaptchaSvg(text: string, width: number, height: number) {
  const safeText = escapeSvg(text || "");
  const lines = Array.from({ length: 4 }, (_, index) => {
    const y1 = 6 + index * 7;
    const y2 = 10 + index * 6;
    const opacity = 0.22 + index * 0.08;
    return `<line x1="8" y1="${y1}" x2="${width - 8}" y2="${y2}" stroke="#8a5cff" stroke-width="1.2" stroke-opacity="${opacity}" />`;
  }).join("");

  const dots = Array.from({ length: 16 }, (_, index) => {
    const x = 10 + (index * 7) % (width - 22);
    const y = 9 + (index * 11) % (height - 14);
    const r = 0.8 + (index % 3) * 0.3;
    const fill = index % 2 === 0 ? "#3b82f6" : "#f472b6";
    return `<circle cx="${x}" cy="${y}" r="${r}" fill="${fill}" fill-opacity="0.35" />`;
  }).join("");

  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
      <rect width="100%" height="100%" rx="4" fill="transparent" />
      ${lines}
      ${dots}
      <text
        x="50%"
        y="58%"
        text-anchor="middle"
        font-family="Courier New, monospace"
        font-size="20"
        font-weight="700"
        letter-spacing="2"
        fill="#9b59b6"
        transform="rotate(-5 ${width / 2} ${height / 2})"
      >${safeText}</text>
    </svg>
  `;

  return `data:image/svg+xml;base64,${window.btoa(unescape(encodeURIComponent(svg)))}`;
}

/**
 * 绘制图形验证码
 * @param width - 图形宽度
 * @param height - 图形高度
 */
export const useImageVerify = (width = 120, height = 40) => {
  const domRef = ref<HTMLImageElement>();
  const imgCode = ref("");
  const captchaKey = ref("");
  const loading = ref(false);

  function setImgCode(code: string) {
    imgCode.value = code;
  }

  function setImgSrc(code: string) {
    if (domRef.value) {
      domRef.value.src = buildCaptchaSvg(code, width, height);
    }
  }

  async function fetchCaptcha() {
    if (loading.value) {
      return;
    }

    loading.value = true;

    try {
      const res = await getCaptcha();
      if (res?.code === 200) {
        captchaKey.value = res.data.captcha_id;
        imgCode.value = res.data.captcha;
        setImgSrc(res.data.captcha);
        return;
      }

      message(res?.message || "获取验证码失败", { type: "error" });
    } catch (error) {
      console.error("获取验证码失败:", error);
      message("获取验证码失败", { type: "error" });
    } finally {
      loading.value = false;
    }
  }

  const getImgCode = debounce(fetchCaptcha, 800, true);

  onMounted(() => {
    fetchCaptcha();
  });

  return {
    domRef,
    imgCode,
    captchaKey,
    loading,
    setImgCode,
    getImgCode
  };
};
