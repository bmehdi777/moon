import { Suspense, lazy } from "react";
import type { ClassKey } from "keycloakify/login";
import type { KcContext } from "./KcContext";
import { useI18n } from "./i18n";
import DefaultPage from "keycloakify/login/DefaultPage";
import Template from "keycloakify/login/Template";
import "./main.css";

const UserProfileFormFields = lazy(
  () => import("keycloakify/login/UserProfileFormFields")
);

const doMakeUserConfirmPassword = true;

export default function KcPage(props: { kcContext: KcContext }) {
  const { kcContext } = props;

  const { i18n } = useI18n({ kcContext });

  return (
    <Suspense>
      {(() => {
        switch (kcContext.pageId) {
          default:
            return (
              <DefaultPage
                kcContext={kcContext}
                i18n={i18n}
                classes={classes}
                Template={Template}
                doUseDefaultCss={true}
                UserProfileFormFields={UserProfileFormFields}
                doMakeUserConfirmPassword={doMakeUserConfirmPassword}
              />
            );
        }
      })()}
    </Suspense>
  );
}

const classes = {
  kcLoginClass: "",
  kcHeaderClass: "",
  kcContentClass: "",

  kcFormCardClass: "",
  kcFormHeaderClass: "",

  kcButtonClass: "",
  kcButtonPrimaryClass: "",
  kcButtonBlockClass: "",
  kcButtonLargeClass: "",
	kcFormPasswordVisibilityButtonClass: "",

  kcFormGroupClass: "",
  kcFormSettingClass: "",
  kcFormOptionsWrapperClass: "",

  kcSignUpClass: "",
  kcInfoAreaWrapperClass: "",
  kcFormSocialAccountListClass: "",
  kcFormSocialAccountListButtonClass: "",

	kcAlertClass: "",
	kcFeedbackInfoIcon: "",
	kcFeedbackErrorIcon: "",
} satisfies { [key in ClassKey]?: string };
