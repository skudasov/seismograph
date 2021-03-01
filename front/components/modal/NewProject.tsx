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
import { useMutation, useQueryClient } from 'react-query'
import { CreateProjectDTO } from '../../api/dto'
import { apiCreateProject } from '../../api/projects'

interface NewProjectModalProps {
  popup: boolean
  setPopup: (flag: boolean) => void
}

const NewProjectModal = (props: NewProjectModalProps) => {
  const { popup, setPopup } = props

  const queryClient = useQueryClient()

  const mut = useMutation(apiCreateProject, {
    onSuccess: () => {
      queryClient.invalidateQueries('projects')
    },
  })

  const validation = yup.object().shape({
    name: yup.string().required('Name is required'),
    description: yup.string().required('Description is required'),
    repoUrl: yup
      .string()
      .matches(
        /((https?):\/\/)?(www.)?[a-z0-9]+(\.[a-z]{2,}){1,3}(#?\/?[a-zA-Z0-9#]+)*\/?(\?[a-zA-Z0-9-_]+=[a-zA-Z0-9-%]+&?)?$/,
        'Enter correct url!'
      )
      .required('Repo url is required'),
  })

  const formik = useFormik({
    initialValues: {
      name: '',
      description: '',
      repoUrl: '',
    },
    validationSchema: validation,
    onSubmit: (values) => {
      mut.mutate({
        name: values.name,
        description: values.description,
        // eslint-disable-next-line @typescript-eslint/camelcase
        repo_url: values.repoUrl,
      } as CreateProjectDTO)
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
            data-test="projects-modal-name-inp"
          />
          {formik.errors.name && formik.touched.name ? (
            <div>{formik.errors.name}</div>
          ) : null}
        </div>
        <div>
          <TextField
            id="description"
            label="Description"
            value={formik.values.description}
            onInput={formik.handleChange}
            data-test="projects-modal-desc-inp"
          />
          {formik.errors.description && formik.touched.description ? (
            <div>{formik.errors.description}</div>
          ) : null}
        </div>
        <div>
          <TextField
            id="repoUrl"
            label="Repository URL"
            value={formik.values.repoUrl}
            onInput={formik.handleChange}
            data-test="projects-modal-repo-inp"
          />
          {formik.errors.repoUrl && formik.touched.repoUrl ? (
            <div>{formik.errors.repoUrl}</div>
          ) : null}
        </div>
      </DialogContent>
      <DialogActions>
        <Button
          onClick={formik.handleSubmit}
          data-test="projects-modal-create-btn"
        >
          Create
        </Button>
        <Button
          onClick={() => {
            setPopup(false)
            formik.resetForm()
          }}
          data-test="projects-modal-cancel"
        >
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  )
}

export default NewProjectModal
