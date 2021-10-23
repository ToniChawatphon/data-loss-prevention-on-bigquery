import os
import google.cloud.dlp
from datetime import datetime


class Dlp:
    service_account: str

    def __init__(self, service_account: str):
        self.service_account = service_account
        self.dlp_client = None
        self.reponse    = None
        self.operation  = None


    def conn(self):
        print('Connect to DLP')
        try:
            os.environ['GOOGLE_APPLICATION_CREDENTIALS'] = self.service_account
            self.dlp_client = google.cloud.dlp_v2.DlpServiceClient()
            print('- connected to dlp')
            return self.dlp_client
        except Exception as ex:
            raise Exception('could not connect to dlp', ex)


    def inspect_content(self, content: str, project: str):
        # TODO add config into config file

        item   = {"value": content}
        parent = f"projects/{project}"
        # The info types to search for in the content. Required.
        info_types = [{"name": "FIRST_NAME"}, {"name": "LAST_NAME"}]
        # The minimum likelihood to constitute a match. Optional.
        min_likelihood = google.cloud.dlp_v2.Likelihood.LIKELIHOOD_UNSPECIFIED
        # The maximum number of findings to report (0 = server maximum). Optional.
        max_findings = 0

        inspect_config = {
                            "info_types": info_types,
                            "min_likelihood": min_likelihood,
                            "limits": {
                                        "max_findings_per_request": max_findings
                                    },
                        }
        request  = {
                        "parent": parent,
                        "inspect_config": inspect_config,
                        "item": item
                    }

        if not self.dlp_client:
            raise Exception('No dlp connection')

        self.response = self.dlp_client.inspect_content(request=request)
        self._display_response_results()


    def create_dlp_job(self, project: str, dataset_id: str, table_id: str):
        # TODO add config into config file
        
        # The info types to search for in the content. Required.
        info_types = [{"name": "FIRST_NAME"}, {"name": "LAST_NAME"}]
        # The minimum likelihood to constitute a match. Optional.
        min_likelihood = google.cloud.dlp_v2.Likelihood.LIKELIHOOD_UNSPECIFIED
        # The maximum number of findings to report (0 = server maximum). Optional.
        max_findings = 0

        inspect_config = {
                            "info_types": info_types,
                            "min_likelihood": min_likelihood,
                            "limits": {
                                "max_findings_per_request": max_findings
                            },
                        }

        storage_config = {
                            "big_query_options": {
                                "table_reference": {
                                    "project_id": project,
                                    "dataset_id": dataset_id,
                                    "table_id": table_id,
                                }
                            }
                        }

        inspect_job = {
                        "inspect_config": inspect_config,
                        "storage_config": storage_config,
                    }

        # parent = f"projects/{project}/locations/global"
        job_id  = '{}_{}_{}'.format(dataset_id, table_id, datetime.now().strftime('%d%m%Y%H%M'))
        parent  = f"projects/{project}"
        request = {
                        "parent": parent, 
                        "inspect_job": inspect_job, 
                        "job_id":job_id
                    }

        if not self.dlp_client:
            raise Exception('No dlp connection')

        self.operation = self.dlp_client.create_dlp_job(request=request)
        print("Inspection operation started: {}".format(self.operation.name))
        self._display_inspect_result()


    def _display_response_results(self):
        if self.response.result.findings:
            for finding in self.response.result.findings:
                try:
                    print("Quote: {}".format(finding.quote))
                except AttributeError as err:
                    print(err)
                print("Info type: {}".format(finding.info_type.name))
                likelihood = finding.likelihood.name
                print("Likelihood: {}".format(likelihood))
        else:
            print("No findings.")


    def _display_inspect_result(self):
        request = {"name": self.operation.name}
        job     = self.dlp_client.get_dlp_job(request=request)
        if job.inspect_details.result.info_type_stats:
            for finding in job.inspect_details.result.info_type_stats:
                print("Info type: {}; Count: {}".format(finding.info_type.name, finding.count))
            else:
                print("No findings.")
