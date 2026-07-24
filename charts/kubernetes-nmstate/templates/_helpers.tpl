{{- define "kubernetes-nmstate.operatorImage" -}}
{{- .Values.operator.image | default (printf "quay.io/nmstate/kubernetes-nmstate-operator:%s" .Chart.AppVersion) -}}
{{- end -}}

{{- define "kubernetes-nmstate.handlerImage" -}}
{{- .Values.handler.image | default (printf "quay.io/nmstate/kubernetes-nmstate-handler:%s" .Chart.AppVersion) -}}
{{- end -}}

{{- define "kubernetes-nmstate.pluginImage" -}}
{{- .Values.plugin.image | default "quay.io/nmstate/nmstate-console-plugin:release-1.0.0" -}}
{{- end -}}
