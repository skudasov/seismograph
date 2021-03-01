import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  TextField,
} from '@material-ui/core'
import React from 'react'
import { useFormik } from 'formik'
import * as yup from 'yup'
// eslint-disable-next-line @typescript-eslint/ban-ts-ignore
// @ts-ignore
import { v4 as uuidv4 } from 'uuid'
import { useMutation, useQueryClient } from 'react-query'
import { ClusterSpecDTO, CreateClusterDTO, InstanceDTO } from '../../api/dto'
import { apiCreateCluster } from '../../api/clusters'

interface NewClusterModalProps {
  popup: boolean
  setPopup: (flag: boolean) => void
}

const NewClusterModal = (props: NewClusterModalProps) => {
  const { popup, setPopup } = props

  const queryClient = useQueryClient()

  const mut = useMutation(apiCreateCluster, {
    onSuccess: () => {
      queryClient.invalidateQueries('clusters')
    },
  })

  const validation = yup.object().shape({
    name: yup.string().required('Name is required'),
    // eslint-disable-next-line @typescript-eslint/camelcase
    provider_name: yup.string().required('Provider name is required'),
    // eslint-disable-next-line @typescript-eslint/camelcase
    project_id: yup.string().required('Project ID is required'),
  })

  const formik = useFormik({
    initialValues: {
      name: '',
      // eslint-disable-next-line @typescript-eslint/camelcase
      provider_name: '',
      // eslint-disable-next-line @typescript-eslint/camelcase
      project_id: 0,
      // eslint-disable-next-line @typescript-eslint/camelcase
      cluster_spec: {} as ClusterSpecDTO,
    } as CreateClusterDTO,
    validationSchema: validation,
    onSubmit: (values) => {
      mut.mutate({
        name: values.name,
        // eslint-disable-next-line @typescript-eslint/camelcase
        provider_name: values.provider_name,
        // eslint-disable-next-line @typescript-eslint/camelcase
        project_id: Number(values.project_id),
        // eslint-disable-next-line @typescript-eslint/camelcase
        cluster_spec: {
          region: 'us-east-2',
          instances: [
            {
              region: 'us-east-2',
              name: `vm-${uuidv4()}`,
              image: 'ami-0c0415cdff14e2a4a',
              type: 't2.micro',
            } as InstanceDTO,
          ],
        } as ClusterSpecDTO,
      } as CreateClusterDTO)
      setPopup(false)
      formik.resetForm()
    },
  })

  return (
    <Dialog open={popup}>
      <DialogContent>
        <div>
          <TextField
            id="name"
            label="Name"
            value={formik.values.name}
            onInput={formik.handleChange}
            autoFocus
            data-test="clusters-modal-name-inp"
          />
          {formik.errors.name && formik.touched.name ? (
            <div>{formik.errors.name}</div>
          ) : null}
        </div>
        <div>
          <TextField
            id="provider_name"
            label="Provider name"
            value={formik.values.provider_name}
            onInput={formik.handleChange}
            data-test="clusters-modal-provider-name-inp"
          />
          {formik.errors.provider_name && formik.touched.provider_name ? (
            <div>{formik.errors.provider_name}</div>
          ) : null}
        </div>
        <div>
          <TextField
            id="project_id"
            label="Project ID"
            value={formik.values.project_id}
            onInput={formik.handleChange}
            data-test="clusters-modal-project-id-inp"
          />
          {formik.errors.project_id && formik.touched.project_id ? (
            <div>{formik.errors.project_id}</div>
          ) : null}
        </div>
      </DialogContent>
      <DialogActions>
        <Button
          onClick={formik.handleSubmit}
          data-test="clusters-modal-create-btn"
        >
          Create
        </Button>
        <Button
          onClick={() => {
            setPopup(false)
            formik.resetForm()
          }}
          data-test="clusters-modal-cancel"
        >
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  )
}

export default NewClusterModal
