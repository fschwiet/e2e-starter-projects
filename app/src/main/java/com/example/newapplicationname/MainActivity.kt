package com.example.newapplicationname

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.res.stringResource

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            MaterialTheme {
                Surface {
                    Greeting()
                }
            }
        }
    }
}

@Composable
@Suppress("ktlint:standard:function-naming", "FunctionNaming")
private fun Greeting() {
    Text(text = stringResource(id = R.string.hello_world))
}
